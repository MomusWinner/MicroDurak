package cases

import (
	"context"

	"github.com/google/uuid"

	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/props"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/utils"
)

func (uc *AuthUseCase) Register(args props.RegisterReq) (resp props.RegisterResp, err error) {
	user, err := uc.ctx.Connection().AuthRepository().GetByEmail(args.Email)
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	if user != nil {
		err = ErrEmailAlreadyTaken
		return
	}

	rep, err := uc.playersClient.CreatePlayer(context.Background(), &players.CreatePlayerRequest{
		Name: args.Name,
		Age:  int32(args.Age),
	})

	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	if uc.smtp == nil {
		uc.ctx.Logger().Error("Smtp is nil")
		err = ErrInternal
		return
	} else {
		err = uc.smtp.Send(args.Email, args.Name)
		uc.ctx.Logger().Info("Send email")
		if err != nil {
			uc.ctx.Logger().Error(err.Error())
		}
	}

	hashedPassword, err := utils.HashPassword(args.Password)
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	playerId, err := uuid.Parse(rep.Id)
	uc.ctx.Logger().Info(playerId.String())
	err = uc.ctx.Connection().AuthRepository().Add(&models.AuthUser{PlayerId: playerId, Email: args.Email, Password: hashedPassword})

	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	jwt, err := utils.GenerateToken(uc.ctx.Config().GetJwtPrivate(), playerId.String())
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	resp = props.RegisterResp{
		PlayerId: playerId,
		Token:    jwt,
	}
	return
}
