package amqppool

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type ChannelPool struct {
	conn     *amqp.Connection
	mu       sync.Mutex
	free     []*amqp.Channel
	maxConns int
}

func NewChannelPool(conn *amqp.Connection, maxConns int) *ChannelPool {
	if maxConns < 1 || maxConns > 100 {
		panic("Channel pool size should be in range from 1 to 100")
	}

	pool := &ChannelPool{conn: conn, maxConns: maxConns}

	for range maxConns {
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to create channel: %v", err)
			continue
		}
		pool.free = append(pool.free, ch)
	}
	return pool
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.free) > 0 {
		ch := p.free[0]
		p.free = p.free[1:]
		return ch, nil
	}

	return p.conn.Channel()
}

func (p *ChannelPool) Return(ch *amqp.Channel) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.free) < p.maxConns {
		p.free = append(p.free, ch)
	} else {
		ch.Close()
	}
}

func (p *ChannelPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ch := range p.free {
		ch.Close()
	}
	p.free = nil
}
