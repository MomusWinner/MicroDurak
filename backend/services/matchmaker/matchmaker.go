package matchmaker

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MommusWinner/MicroDurak/services/matchmaker/config"
	rc "github.com/MommusWinner/MicroDurak/services/matchmaker/redis"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/types"
	"github.com/redis/go-redis/v9"
)

const increaceRangeAfter = 5
const increaceRangeBy = 100
const groupSize = 2

type Matchmaker struct {
	queueChan    <-chan types.MatchChan
	queue        map[string]types.MatchChan
	config       *config.Config
	playerClient *rc.PlayerClient
}

func New(queueChan <-chan types.MatchChan, config *config.Config, redisClient *redis.Client) *Matchmaker {
	playerClient := rc.NewClient(redisClient)
	queue := make(map[string]types.MatchChan)
	return &Matchmaker{queueChan, queue, config, playerClient}
}

func (m *Matchmaker) Start(
	ctx context.Context,
) error {
	fmt.Printf("Starting matchmaker\n\n")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case matchChan := <-m.queueChan:
				m.queue[matchChan.PlayerId] = matchChan
			}
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cancelling matchmaker")
			return nil
		case <-ticker.C:
			if len(m.queue) == 0 {
				continue
			}
			m.matchmake(ctx)
			fmt.Println("----------")
		}
	}
}

func (m *Matchmaker) matchmake(ctx context.Context) error {
	for playerId, player := range m.queue {
		fmt.Printf("Processing player: %s\n", playerId)

		storedPlayer, err := m.playerClient.GetPlayer(ctx, playerId)
		if err != nil {
			storedPlayer, _ = m.playerClient.AddPlayer(ctx, playerId, player.Rating)
		}
		fmt.Printf("Status: %d\n", storedPlayer.Status)

		switch storedPlayer.Status {
		case rc.StatusSearch:
			err := m.handleSearch(ctx, player)
			if errors.Is(err, types.ErrGroupNotFound) {
				continue
			} else if err != nil {
				return err
			}
		case rc.StatusMoved:
			err := m.handleMoved(ctx, storedPlayer)
			var gidError types.ErrGroupTooSmall
			if errors.As(err, &gidError) {
				continue
			} else if err != nil {
				return err
			}
		}
		fmt.Println("")
	}
	return nil
}

func (m *Matchmaker) handleSearch(ctx context.Context, player types.MatchChan) error {
	scoreRange := int(max((time.Now().Unix()-player.SentTime.Unix())/increaceRangeAfter, 1) * increaceRangeBy)
	low := player.Rating - scoreRange
	high := player.Rating + scoreRange
	fmt.Printf("Low %d; High %d\n", low, high)

	count, err := m.playerClient.CountGroups(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		groups, err := m.playerClient.ListGroupsRange(ctx, low, high)
		if err != nil {
			return err
		}

		if len(groups) > 0 {
			groupId, err := rc.ParseGroupId(groups[0].Member.(string))
			if err != nil {
				return err
			}

			m.playerClient.AddToGroup(ctx, groupId, redis.Z{Score: float64(player.Rating), Member: player.PlayerId})
			fmt.Printf("Found group: %v\n", groups[0])
			return nil
		}
	} else {
		players, err := m.playerClient.ListPlayersRange(ctx, low, high)
		if err != nil {
			return err
		}

		if len(players) <= 1 {
			return types.ErrGroupNotFound
		}
		fmt.Printf("Found players: %v\n", players)

		err = m.playerClient.AddGroup(ctx, players[:min(groupSize, len(players))])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Matchmaker) handleMoved(
	ctx context.Context,
	storedPlayer rc.RedisPlayer,
) error {
	len, err := m.playerClient.GetGroupLen(ctx, storedPlayer.Gid)
	if err != nil {
		return err
	}

	if len < groupSize {
		return types.NewGroupTooSmall(storedPlayer.Gid)
	}

	grouppedPlayers, err := m.playerClient.GetGrouppedPlayers(ctx, storedPlayer.Gid, groupSize)
	if err != nil {
		return err
	}

	err = m.playerClient.RemoveGroup(ctx, storedPlayer.Gid)
	if err != nil {
		return err
	}

	response := types.MatchResponse{
		Status: types.MatchFound,
		RoomId: strconv.Itoa(storedPlayer.Gid),
	}

	for _, grouppedPlayer := range grouppedPlayers {
		m.playerClient.SetPlayerStatus(ctx, grouppedPlayer, rc.StatusEmpty)

		fmt.Printf("Sending to player: %s\n", grouppedPlayer)
		m.queue[grouppedPlayer].ReturnChan <- response
		delete(m.queue, grouppedPlayer)
	}
	return nil
}
