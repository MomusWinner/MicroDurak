package main

import (
	"log"

	"github.com/MommusWinner/MicroDurak/services/game/config"
	"github.com/MommusWinner/MicroDurak/services/game/controller"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	opt, err := redis.ParseURL(conf.RedisURL)
	if err != nil {
		return err
	}

	client := redis.NewClient(opt)
	channel, err := connectToRabbit(conf)

	if err != nil {
		return err
	}

	gameController := controller.NewGameController(conf, channel, client)
	gameController.CreateGame([]string{
		"test5",
		"test7",
	})

	if err != nil {
		log.Fatal(err)
	}

	go gameController.ProcessQueues()
	return nil
}

func connectToRabbit(conf *config.Config) (*amqp.Channel, error) {
	conn, err := amqp.Dial(conf.RabbitmqURL)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, err
}

func main() {
	var forever chan struct{}
	if err := run(); err != nil {
		panic(err)
	}
	<-forever
}
