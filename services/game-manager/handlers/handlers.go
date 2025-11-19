package handlers

import (
	"log"
	"time"

	"github.com/MommusWinner/MicroDurak/lib/jwt"
	"github.com/MommusWinner/MicroDurak/services/game-manager/config"
	"github.com/MommusWinner/MicroDurak/services/game-manager/metrics"
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

	gameId := c.Param("gameId")
	if gameId == "" {
		c.Response().Status = 400
		return nil
	}

	if err != nil {
		return err
	}
	defer ws.Close()

	metrics.PlayersConnected.WithLabelValues(h.Config.PodName, h.Config.Namespace).Inc()
	defer metrics.PlayersConnected.WithLabelValues(h.Config.PodName, h.Config.Namespace).Dec()

	endRead := make(chan bool)
	defer close(endRead) // Закрываем канал при выходе из функции

	// Запускаем горутину для чтения сообщений
	go func() {
		for {
			select {
			case <-endRead:
				return
			default:
				_, msg, err := ws.ReadMessage()
				if err != nil {
					c.Logger().Error(err)
					return
				}
				log.Printf("ReadMessage: %v", string(msg))
				publisher.SendMessageToGame(h.Channel, msg)
				log.Printf("%s\n", msg)
			}
		}
	}()

	// Основной цикл для записи сообщений
	for {
		h.processQueue(gameId, userId, func(message []byte) {
			if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
				c.Logger().Error(err)
			}
		})

		// Небольшая пауза, чтобы не нагружать CPU
		time.Sleep(10 * time.Millisecond)
	}
}

func (h Handler) processQueue(gameId string, userId string, processMessage func([]byte)) {
	queue_name := "game-manager-" + userId + "_" + gameId
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
