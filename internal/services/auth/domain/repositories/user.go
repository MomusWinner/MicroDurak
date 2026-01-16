package repositories

import (
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Add(model *models.User) error
}

type AuthRepository interface {
	Add(model *models.AuthUser) error
	Delete(id uuid.UUID) error
	GetByEmail(email string) (*models.AuthUser, error)
}
