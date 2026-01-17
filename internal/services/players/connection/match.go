package connection

import (
	"context"
	"fmt"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type matchRepo struct {
	queries *database.Queries
	conn    *pgx.Conn
}

func NewMatchRepository(conn *pgx.Conn, queries *database.Queries) *matchRepo {
	return &matchRepo{queries: queries, conn: conn}
}

func (r *matchRepo) Add(ctx context.Context, playerCount int, gameResult models.GameResult) (*models.Match, error) {
	var result string = models.GameResult_name[int32(gameResult)]
	match, err := r.queries.CreateMatchResult(ctx, database.CreateMatchResultParams{PlayerCount: int16(playerCount), GameResult: database.GameResult(result)})
	if err != nil {
		return nil, err
	}
	return &models.Match{Id: match.ID, PlayerCount: int(match.PlayerCount), GameResult: gameResult}, nil
}

func (r *matchRepo) AddPlayerToMatch(ctx context.Context, matchId, playerId uuid.UUID, playerPlace int, ratingChange int32) error {
	_, err := r.queries.AddPlayerPlacement(ctx, database.AddPlayerPlacementParams{
		MatchResultID: matchId,
		PlayerID:      playerId,
		PlayerPlace:   int16(playerPlace),
		RatingChange:  ratingChange,
	})

	return err
}

func (r *matchRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	match, err := r.queries.GetMatchResultById(ctx, id)
	if err != nil {
		return nil, err
	}

	gameResult := models.GameResult_value[string(match.GameResult)]

	return &models.Match{
		Id:          match.ID,
		PlayerCount: int(match.PlayerCount),
		GameResult:  models.GameResult(gameResult),
	}, nil
}

func (r *matchRepo) GetAll(ctx context.Context) ([]models.Match, error) {
	matches, err := r.queries.GetAllMatchResults(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Match, len(matches))
	for i, match := range matches {
		gameResult := models.GameResult_value[string(match.GameResult)]
		result[i] = models.Match{
			Id:          match.ID,
			PlayerCount: int(match.PlayerCount),
			GameResult:  models.GameResult(gameResult),
		}
	}

	return result, nil
}

func (r *matchRepo) GetPlayerPlacementsByMatchId(ctx context.Context, matchId uuid.UUID) ([]models.PlayerPlacementWithDetails, error) {
	placements, err := r.queries.GetPlayerPlacementsByMatchId(ctx, matchId)
	if err != nil {
		return nil, err
	}

	result := make([]models.PlayerPlacementWithDetails, len(placements))
	for i, placement := range placements {
		result[i] = models.PlayerPlacementWithDetails{
			PlayerId:      placement.PlayerID,
			PlayerPlace:   int(placement.PlayerPlace),
			RatingChange:  placement.RatingChange,
			PlayerName:    placement.Name,
			CurrentRating: placement.Rating,
		}
	}

	return result, nil
}

func (r *matchRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context, matchRepo repositories.MatchRepository, userRepo repositories.UserRepository) error) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	dbQueries := r.queries.WithTx(tx)

	txMatchRepo := &matchRepo{queries: dbQueries}
	txUserRepo := NewPlayerRepository(dbQueries, r.conn)
	if err := fn(ctx, txMatchRepo, txUserRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
