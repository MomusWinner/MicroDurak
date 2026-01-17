package connection

import (
	"context"
	"fmt"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/infra"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/repositories"
	"github.com/jackc/pgx/v5"
)

type connection struct {
	conn    *pgx.Conn
	queries *database.Queries

	userRepository  repositories.UserRepository
	matchRepository repositories.MatchRepository
}

func makeConnection(conn *pgx.Conn) *connection {
	queries := database.New(conn)

	return &connection{
		queries:         queries,
		matchRepository: NewMatchRepository(conn, queries),
		userRepository:  NewPlayerRepository(queries),
	}
}

func Make(cfg infra.Config) domain.Connection {
	conn, err := pgx.Connect(context.TODO(), cfg.GetDatabaseURL())

	if err != nil {
		panic(fmt.Sprintf("unable to open database due [%s]", err))
	}

	return makeConnection(conn)
}

func Close(conn domain.Connection) {
	c := conn.(*connection)
	c.conn.Close(context.TODO())
}

func (c *connection) UserRepository() repositories.UserRepository {
	return c.userRepository
}

func (c *connection) MatchRepository() repositories.MatchRepository {
	return c.matchRepository
}
