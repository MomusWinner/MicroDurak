package core

import (
	"github.com/MommusWinner/MicroDurak/internal/services/players/delivery/http"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/cases"
)

type Di struct {
	Ctx           domain.Context
	PlayerUseCase *cases.PlayerUseCase
	MatchUseCase  *cases.MatchUseCase
	PlayerHandler *http.PlayerHandler
}

func NewDi() *Di {
	ctx := InitCtx()

	var (
		playerUseCase = cases.NewPlayersUseCase(ctx)
		playerHandler = http.NewPlayerHandler(playerUseCase)
	)

	return &Di{
		Ctx:           ctx,
		PlayerUseCase: playerUseCase,
		MatchUseCase:  cases.NewMatchUseCase(ctx),
		PlayerHandler: playerHandler,
	}
}
