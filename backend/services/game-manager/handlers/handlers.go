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

func (h Handler) processQueue(userId string, processMessage func([]byte)) {
	queue_name := "game-manager-" + userId
	exchange_name := "game-manager-ex"

	_, err := h.Channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Printf("Declare err: %v", err)
		panic(err)
	}
	err = h.Channel.ExchangeDeclare(
		exchange_name, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Printf("Exchange err: %v", err)
		panic(err)
	}
	msgs, _ := h.Channel.Consume(
		queue_name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	func() {
		for d := range msgs {
			processMessage(d.Body)
		}
	}()
}
