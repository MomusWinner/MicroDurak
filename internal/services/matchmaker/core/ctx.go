package core

import (
	"log/slog"
	"os"

	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/connection"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/infra"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/infra/config"
)

type Ctx struct {
	con    domain.Connection
	cfg    infra.Config
	logger *slog.Logger
}

func (c *Ctx) Config() infra.Config {
	return c.cfg
}

func (c *Ctx) Logger() *slog.Logger {
	return c.logger
}

func (c *Ctx) Connection() domain.Connection {
	return c.con
}

func (c *Ctx) Make() domain.Context {
	return &Ctx{
		con:    c.con,
		logger: c.logger,
		cfg:    c.cfg,
	}
}

func InitCtx() *Ctx {
	cfg := config.Make()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	conn := connection.Make(cfg)

	return &Ctx{
		cfg:    cfg,
		logger: logger,
		con:    conn,
	}
}

func DisposeCtx(ctx *Ctx) {
	connection.Close(ctx.con)
}
