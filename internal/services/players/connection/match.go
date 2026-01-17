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

func (r *matchRepo) AddPlayerToMatch(ctx context.Context, matchId, playerId uuid.UUID, ratingChange int32) error {
	_, err := r.queries.AddPlayerPlacement(ctx, database.AddPlayerPlacementParams{
		MatchResultID: matchId,
		PlayerID:      playerId,
		RatingChange:  ratingChange,
	})

	return err
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
