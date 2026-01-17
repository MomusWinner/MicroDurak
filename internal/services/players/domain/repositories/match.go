package repositories

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
)

type MatchRepository interface {
	Add(ctx context.Context, playerCount int, gameResult models.GameResult) (*models.Match, error)
	AddPlayerToMatch(ctx context.Context, matchId, playerId uuid.UUID, playerPlace int, ratingChange int32) error
	WithTransaction(ctx context.Context, fn func(ctx context.Context, matchRepo MatchRepository, userRepo UserRepository) error) error

	GetById(ctx context.Context, id uuid.UUID) (*models.Match, error)
	GetAll(ctx context.Context) ([]models.Match, error)
	GetPlayerPlacementsByMatchId(ctx context.Context, matchId uuid.UUID) ([]models.PlayerPlacementWithDetails, error)
}
