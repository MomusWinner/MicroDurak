package config

import (
	"github.com/alecthomas/kong"
)

type Config struct {
	JWTPublic  string `help:"Base64 Private key for the jwt"          env:"JWT_PUBLIC" required:"true"`
	Port       string `help:"Port to listen on"                       env:"PORT" default:"8080"`
	RedisURL   string `help:"Redis connection URL"                    env:"REDIS_URL" required:"true"`
	PlayersURL string `help:"URL pointing to the Players Service"  env:"PLAYERS_URL" required:"true"`
	GameURL    string `help:"URL pointing to the Game Service"  env:"GAME_URL" required:"true"`
	LogLevel   string `help:"Log level (debug, info, warn, error)"    env:"LOG_LEVEL" default:"info"`
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
