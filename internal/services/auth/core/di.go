package core

import (
	"fmt"

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
	smtp := MakeSMTP(ctx.Config())

	ctx.Logger().Info(fmt.Sprint(smtp))

	var (
		authUseCase = cases.NewAuthUseCase(ctx, playersClient, smtp)
		authHandler = http.NewAuthHandler(authUseCase)
	)

	return &Di{
		Ctx:         ctx,
		AuthHandler: authHandler,
	}
}
