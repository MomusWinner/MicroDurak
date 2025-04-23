package players

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/players/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerService struct {
	players.UnimplementedPlayersServer
	DBConn    *pgxpool.Pool
	DBQueries *database.Queries
	Config    *config.Config
}

func NewPlayerService(dbConn *pgxpool.Pool, config *config.Config) *PlayerService {
	dbQueries := database.New(dbConn)
	return &PlayerService{DBConn: dbConn, DBQueries: dbQueries, Config: config}
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

func (ps *PlayerService) GetPlayer(ctx context.Context, req *players.GetPlayerRequest) (*players.Player, error) {
	playerId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "player_id is not uuid").Err()
	}

	player, err := ps.DBQueries.GetPlayerById(ctx, pgtype.UUID{Valid: true, Bytes: playerId})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	} else if err != nil {
		return nil, err
	}

	return &players.Player{
		Id:     player.ID.String(),
		Name:   player.Name,
		Age:    int32(player.Age),
		Rating: player.Rating,
	}, nil
}

func playerGameResultToString(gr players.GameResult) string {
	switch gr {
	case players.GameResult_DRAW:
		return "draw"
	case players.GameResult_WIN:
		return "win"
	case players.GameResult_INTERRUPTED:
		return "interrupted"
	default:
		panic("unknown game result")
	}
}

func (ps *PlayerService) CreateMatchResult(ctx context.Context, req *players.CreateMatchResultRequest) (*players.CreateMatchResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMatchResult not implemented")
}
