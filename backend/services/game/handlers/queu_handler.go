package handlers

import (
	"log"

	"github.com/MommusWinner/MicroDurak/services/game/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	Config  *config.Config
	Channel *amqp.Channel
}

func (h Handler) ProcessQueues() {
}

func (h Handler) processQueue(processMessage func(string)) {
	queue_name := "game"
	exchange_name := "gameEx"

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
			processMessage(string(d.Body))
		}
	}()
}
