package cases

import (
	"log"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/props"
)

type HandleMessageUseCase struct {
	ctx     domain.Context
	metrics domain.Metrics
}

func NewHandleMessageUseCase(ctx domain.Context, metrics domain.Metrics) *HandleMessageUseCase {
	return &HandleMessageUseCase{
		ctx:     ctx,
		metrics: metrics,
	}
}

func (uc *HandleMessageUseCase) HandleMessage(args props.HandleMessageReq) (resp props.HandleMessageResp, err error) {
	err = uc.ctx.Messaging().SendMessageToGame(args.Message)
	if err != nil {
		uc.ctx.Logger().Error("Failed to send message to game", "error", err.Error())
		err = ErrInternal
		return
	}

	resp = props.HandleMessageResp{
		Success: true,
	}
	return
}

func (uc *HandleMessageUseCase) ConnectWebSocket(args props.ConnectWebSocketReq) error {
	ws := args.WebSocket

	uc.metrics.IncPlayersConnected(uc.ctx.Config().GetPodName(), uc.ctx.Config().GetNamespace())
	defer uc.metrics.DecPlayersConnected(uc.ctx.Config().GetPodName(), uc.ctx.Config().GetNamespace())

	endRead := make(chan bool)
	defer close(endRead)

	go func() {
		for {
			select {
			case <-endRead:
				return
			default:
				msg, err := ws.ReadMessage()
				if err != nil {
					uc.ctx.Logger().Error("Failed to read message", "error", err.Error())
					return
				}
				log.Printf("ReadMessage: %v", string(msg))

				_, err = uc.HandleMessage(props.HandleMessageReq{
					GameId:  args.GameId,
					UserId:  args.UserId,
					Message: msg,
				})
				if err != nil {
					uc.ctx.Logger().Error("Failed to handle message", "error", err.Error())
					return
				}
				log.Printf("%s\n", msg)
			}
		}
	}()

	for {
		err := uc.ctx.Messaging().ProcessQueue(args.GameId, args.UserId, func(message []byte) {
			if err := ws.WriteMessage(message); err != nil {
				uc.ctx.Logger().Error("Failed to write message", "error", err.Error())
			}
		})
		if err != nil {
			uc.ctx.Logger().Error("Failed to process queue", "error", err.Error())
		}

		time.Sleep(10 * time.Millisecond)
	}
}
