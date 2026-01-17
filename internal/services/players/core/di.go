package core

import (
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/cases"
)

type Di struct {
	Ctx           domain.Context
	PlayerUseCase *cases.PlayerUseCase
	MatchUseCase  *cases.MatchUseCase
}

func NewDi() *Di {
	ctx := InitCtx()

	return &Di{
		Ctx:           ctx,
		PlayerUseCase: cases.NewPlayersUseCase(ctx),
		MatchUseCase:  cases.NewMatchUseCase(ctx),
	}
}
