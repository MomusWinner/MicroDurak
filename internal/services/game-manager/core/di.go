package core

import (
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/delivery/http"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/cases"
)

type Di struct {
	Ctx                  domain.Context
	Handler              *http.GameManagerHandler
	HandleMessageUseCase *cases.HandleMessageUseCase
}

func NewDi() *Di {
	ctx := InitCtx()

	var (
		metrics              = http.NewMetricsAdapter()
		handleMessageUseCase = cases.NewHandleMessageUseCase(ctx, metrics)
		handler              = http.NewGameManagerHandler(ctx, handleMessageUseCase)
	)

	return &Di{
		Ctx:                  ctx,
		Handler:              handler,
		HandleMessageUseCase: handleMessageUseCase,
	}
}
