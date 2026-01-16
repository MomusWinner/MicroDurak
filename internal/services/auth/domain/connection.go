package domain

import "github.com/MommusWinner/MicroDurak/internal/services/auth/domain/repositories"

type Connection interface {
	AuthRepository() repositories.AuthRepository
}
