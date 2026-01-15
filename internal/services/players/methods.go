package players

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/config"
	"github.com/MommusWinner/MicroDurak/internal/services/players/rating"
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
	tx, err := ps.DBConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	dbQueries := ps.DBQueries.WithTx(tx)

	match, err := ps.DBQueries.CreateMatchResult(ctx, database.CreateMatchResultParams{
		PlayerCount: int16(len(req.PlayerPlacements)),
		GameResult:  database.GameResult(playerGameResultToString(req.GameResult)),
	})
	if err != nil {
		return nil, err
	}

	playerScores, err := rating.CalculatePlayerScores(ctx, dbQueries, req.PlayerPlacements)
	if err != nil {
		return nil, err
	}

	playerRatings := make([]*players.PlayerPlacementResponse, len(req.PlayerPlacements))
	for i, player := range req.PlayerPlacements {
		playerId, err := uuid.Parse(player.PlayerId)
		if err != nil {
			return nil, status.New(codes.InvalidArgument, "player_id is not a uuid").Err()
		}
		playerDbId := pgtype.UUID{Bytes: playerId, Valid: true}

		playerRating := playerScores[player.PlayerId].NewRating
		playerRatingChange := int32(playerScores[player.PlayerId].RatingChange)

		_, err = dbQueries.AddPlayerPlacement(ctx, database.AddPlayerPlacementParams{
			MatchResultID: match.ID,
			PlayerID:      playerDbId,
			RatingChange:  playerRating,
		})
		if err != nil {
			return nil, err
		}

		_, err = dbQueries.UpdatePlayerRating(ctx, database.UpdatePlayerRatingParams{
			ID:     playerDbId,
			Rating: playerRating,
		})

		playerRatings[i] = &players.PlayerPlacementResponse{
			PlayerId:           playerId.String(),
			PlayerRating:       playerRating,
			PlayerRatingChange: playerRatingChange,
		}
	}
	response := &players.CreateMatchResultResponse{MatchResultId: match.ID.String(), PlayerRatings: playerRatings}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}
