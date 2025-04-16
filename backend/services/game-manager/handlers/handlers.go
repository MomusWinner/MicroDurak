package handlers

import (
	"fmt"

	"github.com/MommusWinner/MicroDurak/services/game-manager/config"
	"github.com/MommusWinner/MicroDurak/services/game-manager/publisher"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	upgrader = websocket.Upgrader{}
)

type Handler struct {
	Config  *config.Config
	Channel *amqp.Channel
}

func AddRoutes(e *echo.Echo, channel *amqp.Channel, config *config.Config) {
	h := Handler{Config: config, Channel: channel}
	e.GET("/game-manager", h.Connect)
}

func (h Handler) Connect(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	request_id := 0

	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		request_id++
		if err != nil {
			c.Logger().Error(err)
		}

		publisher.SendMessageToGame(h.Channel, []byte(msg))

		fmt.Printf("%s\n", msg)
	}
}
