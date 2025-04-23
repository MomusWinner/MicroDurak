package publisher

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendMessageToGame(channel *amqp.Channel, message []byte) error {
	queue_name := "game-commands"
	exchange_name := "gameEx"

	_, err := channel.QueueDeclare(
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

	err = channel.ExchangeDeclare(
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

	err = channel.QueueBind(
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

	err = channel.PublishWithContext(ctx,
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
