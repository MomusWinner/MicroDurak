package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/infra"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type connection struct {
	pool    *pgxpool.Pool
	queries *database.Queries

	userRepository  repositories.UserRepository
	matchRepository repositories.MatchRepository
}

func makeConnection(pool *pgxpool.Pool) *connection {
	queries := database.New(pool)

	return &connection{
		pool:            pool,
		queries:         queries,
		matchRepository: NewMatchRepository(pool, queries),
		userRepository:  NewPlayerRepository(queries, pool),
	}
}

func Make(cfg infra.Config) domain.Connection {
	pool, err := pgxpool.New(context.Background(), cfg.GetDatabaseURL())

	pool.Config().MaxConns = 20
	pool.Config().MinConns = 5
	pool.Config().MaxConnLifetime = time.Hour
	pool.Config().MaxConnIdleTime = 30 * time.Minute

	if err != nil {
		panic(fmt.Sprintf("unable to open database due [%s]", err))
	}

	return makeConnection(pool)
}

func Close(conn domain.Connection) {
	c := conn.(*connection)
	c.pool.Close()
}

func (c *connection) UserRepository() repositories.UserRepository {
	return c.userRepository
}

func (c *connection) MatchRepository() repositories.MatchRepository {
	return c.matchRepository
}
