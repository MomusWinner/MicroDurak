package repositories

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
)

type MatchRepository interface {
	Add(ctx context.Context, playerCount int, gameResult models.GameResult) (*models.Match, error)
	AddPlayerToMatch(ctx context.Context, matchId, playerId uuid.UUID, ratingChange int32) error
	WithTransaction(ctx context.Context, fn func(ctx context.Context, matchRepo MatchRepository, userRepo UserRepository) error) error
}
