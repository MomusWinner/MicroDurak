package players

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/players/config"
)

type contextKey string

const DBSession contextKey = "dbsession"

type PlayerService struct {
	players.UnimplementedPlayersServer
	DBQueries *database.Queries
	Config    *config.Config
}

func NewPlayerService(dbQueries *database.Queries, config *config.Config) *PlayerService {
	return &PlayerService{DBQueries: dbQueries, Config: config}
}

func (ps *PlayerService) CreatePlayer(ctx context.Context, req *players.CreatePlayerRequest) (*players.CreatePlayerReply, error) {
	id, err := ps.DBQueries.CreatePlayer(ctx, database.CreatePlayerParams{
		Name: req.Name,
		Age:  int16(req.Age),
	})
	if err != nil {
		return nil, err
	}

	return &players.CreatePlayerReply{Id: id.String()}, nil
}
