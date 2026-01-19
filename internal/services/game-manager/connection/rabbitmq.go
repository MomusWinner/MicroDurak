package connection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/infra"
	"github.com/MommusWinner/MicroDurak/lib/amqppool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type messaging struct {
	conn *amqp.Connection
	pool *amqppool.ChannelPool
}

func NewMessaging(conn *amqp.Connection) domain.Messaging {
	return &messaging{
		conn: conn,
		pool: amqppool.NewChannelPool(conn, 20),
	}
}

func (m *messaging) ProcessQueue(gameId string, userId string, processMessage func([]byte)) error {
	queue_name := "game-manager-" + userId + "_" + gameId
	exchange_name := "game-manager-ex"

	channel, err := m.pool.Get()
	defer m.pool.Return(channel)

	if err != nil {
		log.Printf("Declare err: %v", err)
		return err
	}

	_, err = channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Printf("Declare err: %v", err)
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
		log.Printf("Exchange err: %v", err)
		return err
	}

	msgs, err := channel.Consume(
		queue_name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("Consume err: %v", err)
		return err
	}

	for d := range msgs {
		processMessage(d.Body)
	}

	return nil
}

func (m *messaging) SendMessageToGame(message []byte) error {
	channel, err := m.pool.Get()
	if err != nil {
		log.Printf("Failed to get channel from pool: %v", err)
		return err
	}
	defer m.pool.Return(channel)

	queue_name := "game"
	exchange_name := "gameEx"

	_, err = channel.QueueDeclare(
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
		exchange_name, // exchange
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

	return nil
}

func Make(cfg infra.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.GetRabbitmqURL())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to rabbitmq due [%s]", err)
	}

	return conn, nil
}

func Close(conn *amqp.Connection) {
	if conn != nil {
		conn.Close()
	}
}
