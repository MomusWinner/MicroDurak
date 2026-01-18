package cases

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/props"
)

type PlayerUseCase struct {
	ctx domain.Context
}

func NewPlayersUseCase(ctx domain.Context) *PlayerUseCase {
	return &PlayerUseCase{
		ctx: ctx,
	}
}

func (uc *PlayerUseCase) Create(args props.CreatePlayerReq) (resp props.CreatePlayerResp, err error) {
	id, err := uc.ctx.Connection().UserRepository().Add(context.Background(), &models.User{Name: args.Name, Age: args.Age})
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	resp = props.CreatePlayerResp{
		Id: id,
	}
	return
}

func (uc *PlayerUseCase) GetById(args props.GetPlayerByIdReq) (resp props.GetPlayerByIdResp, err error) {
	user, err := uc.ctx.Connection().UserRepository().GetById(context.Background(), args.Id)
	if err != nil {

		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	resp = props.GetPlayerByIdResp{
		Player: user,
	}
	return
}

func (uc *PlayerUseCase) GetAll(args props.GetAllPlayersReq) (resp props.GetAllPlayersResp, err error) {
	players, err := uc.ctx.Connection().UserRepository().GetAll(context.Background())
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		err = ErrInternal
		return
	}

	resp = props.GetAllPlayersResp{
		Players: players,
	}
	return
}
