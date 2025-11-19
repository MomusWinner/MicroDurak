package main

import (
	"context"

	"github.com/MommusWinner/MicroDurak/services/game-manager/config"
	"github.com/MommusWinner/MicroDurak/services/game-manager/handlers"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func run(ctx context.Context, e *echo.Echo) error {
	config, err := config.Load()
	if err != nil {
		return err
	}
	e.Logger.Info(config)

	channel, err := connectToRabbit(config)

	if err != nil {
		return err
	}

	handlers.AddRoutes(e, channel, config)

	return e.Start(":" + config.Port)
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
	e := echo.New()
	ctx := context.Background()
	if err := run(ctx, e); err != nil {
		panic(err)
	}
}
