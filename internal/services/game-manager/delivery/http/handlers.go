package http

import (
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/cases"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/props"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

type GameManagerHandler struct {
	ctx                  domain.Context
	handleMessageUseCase *cases.HandleMessageUseCase
}

func NewGameManagerHandler(ctx domain.Context, handleMessageUseCase *cases.HandleMessageUseCase) *GameManagerHandler {
	return &GameManagerHandler{
		ctx:                  ctx,
		handleMessageUseCase: handleMessageUseCase,
	}
}

func (h *GameManagerHandler) Connect(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	userId, ok := c.Get("playerId").(string)
	if !ok {
		if ws != nil {
			ws.Close()
		}
		return echo.NewHTTPError(401, "Unauthorized")
	}

	gameId := c.Param("gameId")
	if gameId == "" {
		c.Response().Status = 400
		if ws != nil {
			ws.Close()
		}
		return nil
	}

	if err != nil {
		if ws != nil {
			ws.Close()
		}
		return err
	}
	defer ws.Close()

	wsAdapter := NewWebSocketAdapter(ws)

	err = h.handleMessageUseCase.ConnectWebSocket(props.ConnectWebSocketReq{
		GameId:    gameId,
		UserId:    userId,
		WebSocket: wsAdapter,
	})

	return err
}
