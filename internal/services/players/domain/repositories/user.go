package repositories

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Add(ctx context.Context, model *models.User) (uuid.UUID, error)
	UpdatePlayerRating(ctx context.Context, playerID uuid.UUID, newRating int) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
}
