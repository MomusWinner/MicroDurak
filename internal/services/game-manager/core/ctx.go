package core

import (
	"log/slog"
	"os"

	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/connection"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/infra"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/infra/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Ctx struct {
	cfg       infra.Config
	logger    *slog.Logger
	messaging domain.Messaging
	conn      *amqp.Connection
}

func (c *Ctx) Config() infra.Config {
	return c.cfg
}

func (c *Ctx) Logger() *slog.Logger {
	return c.logger
}

func (c *Ctx) Messaging() domain.Messaging {
	return c.messaging
}

func (c *Ctx) Make() domain.Context {
	return &Ctx{
		cfg:       c.cfg,
		logger:    c.logger,
		messaging: c.messaging,
		conn:      c.conn,
	}
}

func InitCtx() *Ctx {
	cfg := config.Make()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	conn, err := connection.Make(cfg)
	if err != nil {
		panic(err)
	}

	messaging := connection.NewMessaging(conn)

	return &Ctx{
		cfg:       cfg,
		logger:    logger,
		messaging: messaging,
		conn:      conn,
	}
}

func DisposeCtx(ctx *Ctx) {
	connection.Close(ctx.conn)
}
