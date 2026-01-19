package cases

import (
	"context"
	"errors"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/contracts/game/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/types"
	"github.com/redis/go-redis/v9"
)

const increaceRangeAfter = 5
const increaceRangeBy = 100
const groupSize = 2

type MatchmakerUseCase struct {
	ctx        domain.Context
	queueChan  <-chan types.MatchChan
	cancelChan <-chan types.MatchCancel
	queue      map[string]types.MatchChan
	gameClient game.GameClient
}

func NewMatchmakerUseCase(
	ctx domain.Context,
	queueChan <-chan types.MatchChan,
	cancelChan <-chan types.MatchCancel,
	gameGRPCClient game.GameClient,
) *MatchmakerUseCase {
	queue := make(map[string]types.MatchChan)
	return &MatchmakerUseCase{
		ctx:        ctx,
		queueChan:  queueChan,
		cancelChan: cancelChan,
		queue:      queue,
		gameClient: gameGRPCClient,
	}
}

func (uc *MatchmakerUseCase) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case matchChan := <-uc.queueChan:
				uc.queue[matchChan.PlayerId] = matchChan
			case cancelRequest := <-uc.cancelChan:
				err := uc.ctx.Connection().MatchmakerRepository().RemovePlayer(ctx, cancelRequest.PlayerId)
				if err != nil {
					panic(err)
				}
				delete(uc.queue, cancelRequest.PlayerId)
			}
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if len(uc.queue) == 0 {
				continue
			}
			err := uc.matchmake(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (uc *MatchmakerUseCase) matchmake(ctx context.Context) error {
	repo := uc.ctx.Connection().MatchmakerRepository()

	for playerId, player := range uc.queue {
		storedPlayer, err := repo.GetPlayer(ctx, playerId)
		if err != nil {
			storedPlayer, _ = repo.AddPlayer(ctx, playerId, player.Rating)
		}

		switch storedPlayer.Status {
		case models.StatusSearch:
			player.ReturnChan <- types.MatchResponse{
				Status: types.MatchPending,
			}

			err := uc.handleSearch(ctx, player)
			if errors.Is(err, ErrGroupNotFound) {
				continue
			} else if err != nil {
				return err
			}
		case models.StatusMoved:
			player.ReturnChan <- types.MatchResponse{
				Status: types.MatchFoundGroup,
			}

			err := uc.handleMoved(ctx, storedPlayer)
			var gidError ErrGroupTooSmall
			if errors.As(err, &gidError) {
				continue
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *MatchmakerUseCase) handleSearch(ctx context.Context, player types.MatchChan) error {
	repo := uc.ctx.Connection().MatchmakerRepository()

	scoreRange := int(max((time.Now().Unix()-player.SentTime.Unix())/increaceRangeAfter, 1) * increaceRangeBy)
	low := player.Rating - scoreRange
	high := player.Rating + scoreRange

	count, err := repo.CountGroups(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		groups, err := repo.ListGroupsRange(ctx, low, high)
		if err != nil {
			return err
		}

		if len(groups) > 0 {
			groupId, err := repo.ParseGroupId(groups[0].Member.(string))
			if err != nil {
				return err
			}

			repo.AddToGroup(ctx, groupId, redis.Z{Score: float64(player.Rating), Member: player.PlayerId})
			return nil
		}
	} else {
		players, err := repo.ListPlayersRange(ctx, low, high)
		if err != nil {
			return err
		}

		if len(players) <= 1 {
			return ErrGroupNotFound
		}

		err = repo.AddGroup(ctx, players[:min(groupSize, len(players))])
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *MatchmakerUseCase) handleMoved(
	ctx context.Context,
	storedPlayer models.RedisPlayer,
) error {
	repo := uc.ctx.Connection().MatchmakerRepository()

	len, err := repo.GetGroupLen(ctx, storedPlayer.Gid)
	if err != nil {
		return err
	}

	if len < groupSize {
		return NewGroupTooSmall(storedPlayer.Gid)
	}

	grouppedPlayers, err := repo.GetGrouppedPlayers(ctx, storedPlayer.Gid, groupSize)
	if err != nil {
		return err
	}

	err = repo.RemoveGroup(ctx, storedPlayer.Gid)
	if err != nil {
		return err
	}

	gameId, err := uc.gameClient.CreateGame(ctx, &game.CreateGameRequest{UserIds: grouppedPlayers})
	if err != nil {
		return err
	}

	response := types.MatchResponse{
		Status: types.MatchCreated,
		RoomId: gameId.GameId,
	}

	for _, grouppedPlayer := range grouppedPlayers {
		repo.SetPlayerStatus(ctx, grouppedPlayer, models.StatusEmpty)

		uc.queue[grouppedPlayer].ReturnChan <- response
		delete(uc.queue, grouppedPlayer)
	}
	return nil
}
