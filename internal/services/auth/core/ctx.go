package core

import (
	"log/slog"
	"os"

	"github.com/MommusWinner/MicroDurak/internal/services/auth/connection"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/infra"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/infra/config"
	"github.com/alecthomas/kong"
)

type Config struct {
	JWTPrivate  string `help:"Base64 Private key for the jwt"       env:"JWT_PRIVATE" required:"true"`
	PlayersURL  string `help:"URL pointing to the Players Service"  env:"PLAYERS_URL" required:"true"`
	Port        string `help:"Port to listen on"                    env:"PORT" default:"8080"`
	DatabaseURL string `help:"Database connection URL"              env:"DATABASE_URL" required:"true"`
	LogLevel    string `help:"Log level (debug, info, warn, error)" env:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	parser, err := kong.New(cfg)
	if err != nil {
		return nil, err
	}

	// Parse command-line flags, environment variables, and config file
	_, err = parser.Parse(nil)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

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

	db := connection.Make(cfg)

	return &Ctx{
		cfg:    cfg,
		logger: logger,
		con:    db,
	}
}

func DisposeCtx(ctx *Ctx) {
	connection.Close(ctx.con)
}
