package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const groupQueueKey = "matchmaking:queue:groups"
const playerQueueKey = "matchmaking:queue:players"

const playerMemFmt = "player:%s"
const groupMemFmt = "group:%d"

const groupsAmountKey = "matchmaking:groups:amount"

const groupKeyFmt = "matchmaking:group:%d"
const groupMembersKey = ":members"

const playerKeyFmt = "matchmaking:player:%s"
const playerStatusKey = ":status"
const playerGroupKey = ":group"

type PlayerStatus = int

type parseError struct {
	name string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("parse failed for %s", e.name)
}

func newParseError(name string) *parseError {
	return &parseError{name}
}

const (
	StatusEmpty = iota
	StatusSearch
	StatusMoved
)

type RedisPlayer struct {
	Status PlayerStatus
	Id     string
	Gid    int
}

type PlayerClient struct {
	client *redis.Client
}

func NewClient(client *redis.Client) *PlayerClient {
	return &PlayerClient{client}
}

func ParseGroupId(groupString string) (int, error) {
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

func ParsePlayerId(playerString string) (string, error) {
	sep := strings.Split(playerString, ":")
	if len(sep) < 2 {
		return "", newParseError("player id")
	}

	return sep[1], nil
}

func (pc *PlayerClient) GetPlayerScore(ctx context.Context, playerId string) (int, error) {
	playerKey := fmt.Sprintf(playerMemFmt, playerId)
	score, err := pc.client.ZScore(ctx, playerQueueKey, playerKey).Result()

	if err != nil {
		return 0, err
	}

	return int(score), nil
}

func (pc *PlayerClient) GetPlayer(ctx context.Context, playerId string) (RedisPlayer, error) {
	playerKey := fmt.Sprintf(playerKeyFmt, playerId)
	statusString, err := pc.client.Get(ctx, playerKey+playerStatusKey).Result()

	var player RedisPlayer

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
	case StatusSearch:
		return player, nil
	case StatusMoved:
		gidString, err := pc.client.Get(ctx, playerKey+playerGroupKey).Result()
		if err != nil {
			return RedisPlayer{}, err
		}

		gid, err := strconv.Atoi(gidString)
		if err != nil {
			return RedisPlayer{}, err
		}

		player.Gid = gid

		return player, nil
	}
	return player, errors.New("unknown player status")
}

func (pc *PlayerClient) SetPlayerStatus(ctx context.Context, playerId string, status PlayerStatus) error {
	playerKey := fmt.Sprintf(playerKeyFmt, playerId)
	err := pc.client.Set(ctx, playerKey+playerStatusKey, status, 24*time.Hour).Err()

	if err != nil {
		return err
	}
	return nil
}

func (pc *PlayerClient) SetPlayerGroup(ctx context.Context, playerId string, group int) error {
	playerKey := fmt.Sprintf(playerKeyFmt, playerId)
	err := pc.client.Set(ctx, playerKey+playerGroupKey, group, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (pc *PlayerClient) AddPlayer(ctx context.Context, playerId string, score int) (RedisPlayer, error) {
	player := RedisPlayer{StatusEmpty, playerId, 0}

	playerKey := fmt.Sprintf(playerMemFmt, playerId)
	err := pc.client.ZAdd(ctx, playerQueueKey, redis.Z{Score: float64(score), Member: playerKey}).Err()
	if err != nil {
		return player, err
	}

	err = pc.SetPlayerStatus(ctx, playerId, StatusSearch)
	if err != nil {
		return player, err
	}

	player.Status = StatusSearch
	return player, nil
}

func (pc *PlayerClient) ListPlayersRange(ctx context.Context, low int, high int) ([]redis.Z, error) {
	lows := fmt.Sprint(low)
	highs := fmt.Sprint(high)

	player, err := pc.client.ZRangeByScoreWithScores(ctx, playerQueueKey, &redis.ZRangeBy{Min: lows, Max: highs}).Result()
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (pc *PlayerClient) AddGroup(ctx context.Context, players []redis.Z) error {
	members := make([]any, 0, len(players))
	for _, player := range players {
		members = append(members, player.Member)
	}

	err := pc.client.ZRem(ctx, playerQueueKey, members...).Err()
	if err != nil {
		return nil
	}

	scoreSum := 0
	for _, player := range players {
		scoreSum += int(player.Score)
	}

	scoreAvg := scoreSum / len(players)

	count, err := pc.client.Incr(ctx, groupsAmountKey).Result()
	if err != nil {
		return nil
	}

	playerIds := make([]string, len(players))
	for i, player := range players {
		id, err := ParsePlayerId(player.Member.(string))
		if err != nil {
			return err
		}

		playerIds[i] = id

		pc.SetPlayerStatus(ctx, id, StatusMoved)
		pc.SetPlayerGroup(ctx, id, int(count))
	}

	groupKey := fmt.Sprintf(groupKeyFmt, count)
	err = pc.client.RPush(ctx, groupKey+groupMembersKey, playerIds).Err()
	if err != nil {
		panic(err)
	}

	queueGroupKey := fmt.Sprintf(groupMemFmt, count)
	err = pc.client.ZAdd(ctx, groupQueueKey, redis.Z{Score: float64(scoreAvg), Member: queueGroupKey}).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (pc *PlayerClient) CountGroups(ctx context.Context) (int, error) {
	count, err := pc.client.ZCard(ctx, groupQueueKey).Result()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (pc *PlayerClient) ListGroupsRange(ctx context.Context, low int, high int) ([]redis.Z, error) {
	lows := fmt.Sprint(low)
	highs := fmt.Sprint(high)
	player, err := pc.client.ZRangeByScoreWithScores(ctx, groupQueueKey, &redis.ZRangeBy{Min: lows, Max: highs}).Result()
	if err != nil {
		return nil, err
	}

	return player, err
}

func (pc *PlayerClient) AddToGroup(ctx context.Context, groupId int, player redis.Z) error {
	playerId := player.Member.(string)

	groupKey := fmt.Sprintf(groupKeyFmt, groupId)
	groupQueueKey := fmt.Sprintf(groupMemFmt, groupId)

	pc.SetPlayerStatus(ctx, playerId, StatusMoved)
	pc.SetPlayerGroup(ctx, playerId, groupId)

	len, err := pc.GetGroupLen(ctx, groupId)
	if err != nil {
		return err
	}

	oldScore, err := pc.GetGroupScore(ctx, groupId)
	if err != nil {
		return err
	}

	newScore := (oldScore*len+int(player.Score))/len + 1

	err = pc.client.ZIncrBy(ctx, groupQueueKey, float64(newScore+oldScore), groupQueueKey).Err()
	if err != nil {
		return err
	}

	err = pc.client.RPush(ctx, groupKey+groupMembersKey, playerId).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (pc *PlayerClient) GetGroupLen(ctx context.Context, groupId int) (int, error) {
	groupKey := fmt.Sprintf(groupKeyFmt, groupId)

	len, err := pc.client.LLen(ctx, groupKey+groupMembersKey).Result()
	if err != nil {
		return 0, err
	}

	return int(len), err
}

func (pc *PlayerClient) GetGrouppedPlayers(ctx context.Context, groupId int, amount int) ([]string, error) {
	groupKey := fmt.Sprintf(groupKeyFmt, groupId)
	players, err := pc.client.LRange(ctx, groupKey+groupMembersKey, 0, int64(amount)).Result()
	if err != nil {
		return nil, err
	}

	return players, err
}

func (pc *PlayerClient) RemoveGroup(ctx context.Context, groupId int) error {
	groupKey := fmt.Sprintf(groupMemFmt, groupId)
	err := pc.client.ZRem(ctx, groupQueueKey, groupKey).Err()
	if err != nil {
		return err
	}
	return nil
}

func (pc *PlayerClient) GetGroupScore(ctx context.Context, groupId int) (int, error) {
	groupKey := fmt.Sprintf(groupMemFmt, groupId)
	score, err := pc.client.ZScore(ctx, groupKey, groupKey).Result()
	if err != nil {
		return 0, err
	}

	return int(score), err
}
