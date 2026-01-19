package core

import (
	"github.com/MommusWinner/MicroDurak/internal/contracts/game/v1"
	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/delivery/http"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/cases"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/types"
)

type Di struct {
	Ctx               domain.Context
	MatchmakerUseCase *cases.MatchmakerUseCase
	Handler           *http.Handler
	QueueChan         chan types.MatchChan
	CancelChan        chan types.MatchCancel
	PlayersClient     players.PlayersClient
	GameClient        game.GameClient
}

func NewDi() *Di {
	ctx := InitCtx()
	playersClient := MakePlayersClient(ctx.Config())
	gameClient := MakeGameClient(ctx.Config())

	queueChan := make(chan types.MatchChan)
	cancelChan := make(chan types.MatchCancel)

	matchmakerUseCase := cases.NewMatchmakerUseCase(
		ctx,
		queueChan,
		cancelChan,
		gameClient,
	)

	handler := http.NewHandler(
		queueChan,
		cancelChan,
		ctx,
		playersClient,
	)

	return &Di{
		Ctx:               ctx,
		MatchmakerUseCase: matchmakerUseCase,
		Handler:           handler,
		QueueChan:         queueChan,
		CancelChan:        cancelChan,
		PlayersClient:     playersClient,
		GameClient:        gameClient,
	}
}
