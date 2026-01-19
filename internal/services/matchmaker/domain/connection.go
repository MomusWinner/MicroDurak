package domain

import "github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/repositories"

type Connection interface {
	MatchmakerRepository() repositories.MatchmakerRepository
}
