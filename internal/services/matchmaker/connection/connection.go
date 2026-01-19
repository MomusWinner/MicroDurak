package connection

import (
	"fmt"

	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/infra"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/repositories"
	"github.com/redis/go-redis/v9"
)

type connection struct {
	client                *redis.Client
	matchmakerRepository repositories.MatchmakerRepository
}

func makeConnection(client *redis.Client) *connection {
	return &connection{
		client:                client,
		matchmakerRepository: NewMatchmakerRepository(client),
	}
}

func Make(cfg infra.Config) domain.Connection {
	opt, err := redis.ParseURL(cfg.GetRedisURL())
	if err != nil {
		panic(fmt.Sprintf("unable to parse redis URL due [%s]", err))
	}
	
	client := redis.NewClient(opt)
	return makeConnection(client)
}

func Close(conn domain.Connection) {
	c := conn.(*connection)
	c.client.Close()
}

func (c *connection) MatchmakerRepository() repositories.MatchmakerRepository {
	return c.matchmakerRepository
}
