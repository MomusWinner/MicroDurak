package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MommusWinner/MicroDurak/services/game/config"
	"github.com/MommusWinner/MicroDurak/services/game/core"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type GameController struct {
	Config  *config.Config
	Channel *amqp.Channel
	Redis   *redis.Client
}

func NewGameController(
	conf *config.Config,
	channel *amqp.Channel,
	redis *redis.Client,
) GameController {
	return GameController{
		Config:  conf,
		Channel: channel,
		Redis:   redis,
	}
}

func (gc GameController) CreateGame(userIds []string) {
	core.CreateNewGameAndSaveInRedis(gc.Redis, userIds)
}

func (gc GameController) LoadGame(gameId string) (*core.Game, error) {
	return core.LoadGame(gc.Redis, gameId)
}

func (gc GameController) ProcessQueues() {
	gc.processQueue(func(message []byte) {
		var command core.Command
		json.Unmarshal(message, &command)
		game, err := core.LoadGame(gc.Redis, command.GameId)
		if err != nil {
			log.Fatalf("Room with id %s does not exist", command.GameId)
			return
		}

		messageByUser, err := game.HandleMessage(message)
		if err != nil {
			log.Fatalf("Error occurred while processing a messag. Room id: %s", command.GameId)
			return
		}

		for userId, userMessage := range messageByUser {
			gc.SendMessageToGameManager(userId, userMessage)
		}
	})
}

func (gc GameController) processQueue(processMessage func([]byte)) {
	queue_name := "game"
	exchange_name := "gameEx"

	_, err := gc.Channel.QueueDeclare(
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
	err = gc.Channel.ExchangeDeclare(
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
	msgs, _ := gc.Channel.Consume(
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

func (gc GameController) SendMessageToGameManager(userId string, message []byte) error {
	queue_name := "game-manager-" + userId
	exchange_name := "game-manager-ex"

	_, err := gc.Channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Print(err)
		return err
	}

	err = gc.Channel.ExchangeDeclare(
		exchange_name, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Print(err)
		return err
	}

	err = gc.Channel.QueueBind(
		queue_name,    // queue name
		queue_name,    // routing key
		exchange_name, // exchange
		false,
		nil,
	)

	if err != nil {
		log.Print(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = gc.Channel.PublishWithContext(ctx,
		exchange_name, // exChannelange
		queue_name,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		log.Print(err)
		return err
	}

	return err
}
