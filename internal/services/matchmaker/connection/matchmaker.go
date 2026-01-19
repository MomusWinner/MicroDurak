package connection

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/repositories"
	"github.com/redis/go-redis/v9"
)

const groupQueueKey = "matchmaking:queue:groups"
const playerQueueKey = "matchmaking:queue:players"

const groupMemFmt = "group:%d"

const groupsAmountKey = "matchmaking:groups:amount"

const groupKeyFmt = "matchmaking:group:%d"
const groupMembersKey = ":members"

const playerKeyFmt = "matchmaking:player:%s"
const playerStatusKey = ":status"
const playerGroupKey = ":group"

type parseError struct {
	name string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("parse failed for %s", e.name)
}

func newParseError(name string) *parseError {
	return &parseError{name}
}

type matchmakerRepository struct {
	client *redis.Client
}

func NewMatchmakerRepository(client *redis.Client) repositories.MatchmakerRepository {
	return &matchmakerRepository{client: client}
}

func (r *matchmakerRepository) ParseGroupId(groupString string) (int, error) {
	sep := strings.Split(groupString, ":")
	if len(sep) < 2 {
		return 0, newParseError("group id")
	}

	id, err := strconv.Atoi(sep[1])
	if err != nil {
		return 0, newParseError("group id")
	}

	return id, nil
}

func (r *matchmakerRepository) GetPlayerScore(ctx context.Context, playerId string) (int, error) {
	score, err := r.client.ZScore(ctx, playerQueueKey, playerId).Result()

	if err != nil {
		return 0, err
	}

	return int(score), nil
}

func (r *matchmakerRepository) GetPlayer(ctx context.Context, playerId string) (models.RedisPlayer, error) {
	playerKey := fmt.Sprintf(playerKeyFmt, playerId)
	statusString, err := r.client.Get(ctx, playerKey+playerStatusKey).Result()

	var player models.RedisPlayer

	if err != nil {
		return player, err
	}

	status, err := strconv.Atoi(statusString)

	if err != nil {
		return player, err
	}

	player.Status = status
	player.Id = playerId

	switch status {
	case models.StatusSearch:
		return player, nil
	case models.StatusMoved:
		gidString, err := r.client.Get(ctx, playerKey+playerGroupKey).Result()
		if err != nil {
			return models.RedisPlayer{}, err
		}

		gid, err := strconv.Atoi(gidString)
		if err != nil {
			return models.RedisPlayer{}, err
		}

		player.Gid = gid

		return player, nil
	}
	return player, errors.New("unknown player status")
}

func (r *matchmakerRepository) SetPlayerStatus(ctx context.Context, playerId string, status models.PlayerStatus) error {
	playerKey := fmt.Sprintf(playerKeyFmt, playerId)
	err := r.client.Set(ctx, playerKey+playerStatusKey, status, 24*time.Hour).Err()

	if err != nil {
		return err
	}
	return nil
}

func (r *matchmakerRepository) SetPlayerGroup(ctx context.Context, playerId string, group int) error {
	playerKey := fmt.Sprintf(playerKeyFmt, playerId)
	err := r.client.Set(ctx, playerKey+playerGroupKey, group, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *matchmakerRepository) AddPlayer(ctx context.Context, playerId string, score int) (models.RedisPlayer, error) {
	player := models.RedisPlayer{Status: models.StatusEmpty, Id: playerId, Gid: 0}

	err := r.client.ZAdd(ctx, playerQueueKey, redis.Z{Score: float64(score), Member: playerId}).Err()
	if err != nil {
		return player, err
	}

	err = r.SetPlayerStatus(ctx, playerId, models.StatusSearch)
	if err != nil {
		return player, err
	}

	player.Status = models.StatusSearch
	return player, nil
}

func (r *matchmakerRepository) ListPlayersRange(ctx context.Context, low int, high int) ([]redis.Z, error) {
	lows := fmt.Sprint(low)
	highs := fmt.Sprint(high)

	player, err := r.client.ZRangeByScoreWithScores(ctx, playerQueueKey, &redis.ZRangeBy{Min: lows, Max: highs}).Result()
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (r *matchmakerRepository) RemovePlayer(ctx context.Context, playerId string) error {
	player, err := r.GetPlayer(ctx, playerId)
	if errors.Is(err, redis.Nil) {
		return nil
	} else if err != nil {
		return err
	}

	switch player.Status {
	case models.StatusEmpty:
		return nil
	case models.StatusSearch:
		err := r.SetPlayerStatus(ctx, playerId, models.StatusEmpty)
		if err != nil {
			return err
		}

		err = r.client.ZRem(ctx, playerQueueKey, playerId).Err()
		if err != nil {
			return err
		}
		return nil
	case models.StatusMoved:
		err := r.RemoveFromGroup(ctx, player.Gid, playerId)
		if err != nil {
			return err
		}
		r.SetPlayerStatus(ctx, playerId, models.StatusEmpty)
		return nil
	}

	return nil
}

func (r *matchmakerRepository) AddGroup(ctx context.Context, players []redis.Z) error {
	members := make([]any, 0, len(players))
	for _, player := range players {
		members = append(members, player.Member)
	}

	err := r.client.ZRem(ctx, playerQueueKey, members...).Err()
	if err != nil {
		return err
	}

	scoreSum := 0
	for _, player := range players {
		scoreSum += int(player.Score)
	}

	scoreAvg := scoreSum / len(players)

	count, err := r.client.Incr(ctx, groupsAmountKey).Result()
	if err != nil {
		return err
	}

	playerIds := make([]string, len(players))
	for i, player := range players {
		id := player.Member.(string)

		playerIds[i] = id

		r.SetPlayerStatus(ctx, id, models.StatusMoved)
		r.SetPlayerGroup(ctx, id, int(count))
	}

	groupKey := fmt.Sprintf(groupKeyFmt, count)
	err = r.client.SAdd(ctx, groupKey+groupMembersKey, playerIds).Err()
	if err != nil {
		return err
	}

	queueGroupKey := fmt.Sprintf(groupMemFmt, count)
	err = r.client.ZAdd(ctx, groupQueueKey, redis.Z{Score: float64(scoreAvg), Member: queueGroupKey}).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *matchmakerRepository) CountGroups(ctx context.Context) (int, error) {
	count, err := r.client.ZCard(ctx, groupQueueKey).Result()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *matchmakerRepository) ListGroupsRange(ctx context.Context, low int, high int) ([]redis.Z, error) {
	lows := fmt.Sprint(low)
	highs := fmt.Sprint(high)
	player, err := r.client.ZRangeByScoreWithScores(ctx, groupQueueKey, &redis.ZRangeBy{Min: lows, Max: highs}).Result()
	if err != nil {
		return nil, err
	}

	return player, err
}

func (r *matchmakerRepository) AddToGroup(ctx context.Context, groupId int, player redis.Z) error {
	playerId := player.Member.(string)

	groupKey := fmt.Sprintf(groupKeyFmt, groupId)
	groupQueueKey := fmt.Sprintf(groupMemFmt, groupId)

	r.SetPlayerStatus(ctx, playerId, models.StatusMoved)
	r.SetPlayerGroup(ctx, playerId, groupId)

	len, err := r.GetGroupLen(ctx, groupId)
	if err != nil {
		return err
	}

	oldScore, err := r.GetGroupScore(ctx, groupId)
	if err != nil {
		return err
	}

	newScore := (oldScore*len+int(player.Score))/len + 1

	err = r.client.ZIncrBy(ctx, groupQueueKey, float64(newScore+oldScore), groupQueueKey).Err()
	if err != nil {
		return err
	}

	err = r.client.SAdd(ctx, groupKey+groupMembersKey, playerId).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (r *matchmakerRepository) GetGroupLen(ctx context.Context, groupId int) (int, error) {
	groupKey := fmt.Sprintf(groupKeyFmt, groupId)

	len, err := r.client.SCard(ctx, groupKey+groupMembersKey).Result()
	if err != nil {
		return 0, err
	}

	return int(len), err
}

func (r *matchmakerRepository) GetGrouppedPlayers(ctx context.Context, groupId int, amount int) ([]string, error) {
	groupKey := fmt.Sprintf(groupKeyFmt, groupId)
	players, err := r.client.SMembers(ctx, groupKey+groupMembersKey).Result()
	if err != nil {
		return nil, err
	}

	return players, err
}

func (r *matchmakerRepository) RemoveFromGroup(ctx context.Context, groupId int, playerId string) error {
	err := r.client.SRem(ctx, groupQueueKey, playerId).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *matchmakerRepository) RemoveGroup(ctx context.Context, groupId int) error {
	groupMemKey := fmt.Sprintf(groupMemFmt, groupId)
	err := r.client.ZRem(ctx, groupQueueKey, groupMemKey).Err()
	if err != nil {
		return err
	}

	groupKey := fmt.Sprintf(groupKeyFmt, groupId)
	err = r.client.Del(ctx, groupKey+groupMembersKey).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *matchmakerRepository) GetGroupScore(ctx context.Context, groupId int) (int, error) {
	groupKey := fmt.Sprintf(groupMemFmt, groupId)
	score, err := r.client.ZScore(ctx, groupKey, groupKey).Result()
	if err != nil {
		return 0, err
	}

	return int(score), err
}
