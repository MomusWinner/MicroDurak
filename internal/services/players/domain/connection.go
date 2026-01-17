package domain

import "github.com/MommusWinner/MicroDurak/internal/services/players/domain/repositories"

type Connection interface {
	UserRepository() repositories.UserRepository
	MatchRepository() repositories.MatchRepository
}
