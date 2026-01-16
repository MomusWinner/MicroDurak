package config

import (
	"github.com/alecthomas/kong"
	"log"
)

import ()

type Config struct {
	JWTPrivate  string `help:"Base64 Private key for the jwt"       env:"JWT_PRIVATE" required:"true"`
	PlayersURL  string `help:"URL pointing to the Players Service"  env:"PLAYERS_URL" required:"true"`
	Port        string `help:"Port to listen on"                    env:"PORT" default:"8080"`
	DatabaseURL string `help:"Database connection URL"              env:"DATABASE_URL" required:"true"`
	LogLevel    string `help:"Log level (debug, info, warn, error)" env:"LOG_LEVEL" default:"info"`
}

// func Load() (*Config, error) {
// 	cfg := &Config{}
// 	parser, err := kong.New(cfg)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Parse command-line flags, environment variables, and config file
// 	_, err = parser.Parse(nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return cfg, nil
// }

func Make() *Config {
	cfg := &Config{}
	parser, err := kong.New(cfg)
	if err != nil {
		log.Panic(err)
	}

	// Parse command-line flags, environment variables, and config file
	_, err = parser.Parse(nil)
	if err != nil {
		log.Panic(err)
	}
	return cfg
}

func (s *Config) GetJwtPrivate() string {
	return s.JWTPrivate
}

func (s *Config) GetPlayersURL() string {
	return s.PlayersURL
}

func (s *Config) GetPort() string {
	return s.Port
}
func (s *Config) GetDatabaseURL() string {
	return s.DatabaseURL
}

func (s *Config) GetLogLevel() string {
	return s.LogLevel
}
