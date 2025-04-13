package config

import (
	"github.com/alecthomas/kong"
)

type Config struct {
	JWTPrivate  string `help:"Base64 Private key for the jwt" env:"JWT_PRIVATE" required:"true"`
	Port        string `help:"Port to listen on" env:"PORT" default:"8080"`
	DatabaseURL string `help:"Database connection URL" env:"DATABASE_URL" required:"true"`
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
