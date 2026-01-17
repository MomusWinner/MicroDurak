package core

import (
	"github.com/MommusWinner/MicroDurak/internal/services/auth/delivery/http"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/cases"
)

type Di struct {
	Ctx         domain.Context
	AuthHandler *http.AuthHandler
}

func NewDi() *Di {
	ctx := InitCtx()
	playersClient := MakePlayersClient(ctx.Config())

	var (
		authUseCase = cases.NewAuthUseCase(ctx, playersClient)
		authHandler = http.NewAuthHandler(authUseCase)
	)

	return &Di{
		Ctx:         ctx,
		AuthHandler: authHandler,
	}
}
