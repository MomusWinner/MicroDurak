package cases

import (
	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/props"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/utils"
)

type AuthUseCase struct {
	ctx           domain.Context
	playersClient players.PlayersClient
	smtp          domain.SMTP
}

func NewAuthUseCase(ctx domain.Context, playersClient players.PlayersClient, smtp domain.SMTP) *AuthUseCase {
	return &AuthUseCase{
		ctx:           ctx,
		playersClient: playersClient,
		smtp:          smtp,
	}
}

func (uc *AuthUseCase) Login(args props.LoginReq) (resp props.LoginResp, err error) {
	user, err := uc.ctx.Connection().AuthRepository().GetByEmail(args.Email)

	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	if user == nil {
		err = ErrLoginFailed
		return
	}

	if !utils.CheckPasswordHash(args.Password, user.Password) {
		err = ErrLoginFailed
		return
	}

	jwt, err := utils.GenerateToken(uc.ctx.Config().GetJwtPrivate(), user.PlayerId.String())
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	resp = props.LoginResp{
		PlayerId: user.PlayerId.String(),
		Token:    jwt,
	}
	return
}
