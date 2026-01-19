package repositories

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/models"
	"github.com/redis/go-redis/v9"
)

type MatchmakerRepository interface {
	GetPlayer(ctx context.Context, playerId string) (models.RedisPlayer, error)
	AddPlayer(ctx context.Context, playerId string, score int) (models.RedisPlayer, error)
	RemovePlayer(ctx context.Context, playerId string) error
	ListPlayersRange(ctx context.Context, low int, high int) ([]redis.Z, error)
	CountGroups(ctx context.Context) (int, error)
	ListGroupsRange(ctx context.Context, low int, high int) ([]redis.Z, error)
	AddGroup(ctx context.Context, players []redis.Z) error
	AddToGroup(ctx context.Context, groupId int, player redis.Z) error
	GetGroupLen(ctx context.Context, groupId int) (int, error)
	GetGrouppedPlayers(ctx context.Context, groupId int, amount int) ([]string, error)
	RemoveGroup(ctx context.Context, groupId int) error
	SetPlayerStatus(ctx context.Context, playerId string, status models.PlayerStatus) error
	ParseGroupId(groupString string) (int, error)
}
