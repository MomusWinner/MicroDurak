package main

import (
	"log"

	"github.com/MommusWinner/MicroDurak/services/game/config"
	"github.com/MommusWinner/MicroDurak/services/game/core"
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
	// game, err := core.CreateNewGame(client, []string{
	// 	"test1",
	// 	"test2",
	// })
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _ = game

	roomId := "7fa88043-5d53-414d-9849-de0e49f5996b"
	game, err := core.LoadGame(client, roomId)

	if err != nil {
		log.Fatal(err)
	}
	_ = game

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
