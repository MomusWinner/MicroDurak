package handlers

import (
	"log"

	"github.com/MommusWinner/MicroDurak/lib/jwt"
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
	e.GET("/game-manager/:gameId", h.Connect, jwt.AuthMiddleware(config.JWTPublic))
}

func (h Handler) Connect(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	userId, ok := c.Get("playerId").(string)
	if !ok {
		panic("Add auth middleware")
	}

	if err != nil {
		return err
	}
	defer ws.Close()

	endRead := make(chan bool)

	for {
		// Write
		if err != nil {
			c.Logger().Error(err)
		}

		h.processQueue(userId, func(message []byte) {
			ws.WriteMessage(websocket.TextMessage, message)
		})

		go func() {
			for {
				select {
				case <-endRead:
					return
				default:
					_, msg, err := ws.ReadMessage()
					if err != nil {
						c.Logger().Error(err)
					}

					publisher.SendMessageToGame(h.Channel, []byte(msg))
					log.Printf("%s\n", msg)
				}
			}
		}()

		endRead <- true
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
