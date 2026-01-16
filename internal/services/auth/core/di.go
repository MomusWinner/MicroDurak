package core

import (
	v1 "github.com/MommusWinner/MicroDurak/internal/services/auth/api/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/cases"
)

type Di struct {
	Ctx         domain.Context
	AuthHandler *v1.AuthHandler
}

func NewDi() *Di {
	ctx := InitCtx()
	playersClient := MakePlayersClient(ctx.Config())

	var (
		authUseCase = cases.NewAuthUseCase(ctx, playersClient)
		authHandler = v1.NewAuthHandler(authUseCase)
	)

	return &Di{
		Ctx:         ctx,
		AuthHandler: authHandler,
	}
}
