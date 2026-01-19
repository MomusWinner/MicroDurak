package config

import (
	"log"

	"github.com/alecthomas/kong"
)

type Config struct {
	JWTPublic  string `help:"Base64 Private key for the jwt"          env:"JWT_PUBLIC" required:"true"`
	Port       string `help:"Port to listen on"                       env:"PORT" default:"8080"`
	RedisURL   string `help:"Redis connection URL"                    env:"REDIS_URL" required:"true"`
	PlayersURL string `help:"URL pointing to the Players Service"  env:"PLAYERS_URL" required:"true"`
	GameURL    string `help:"URL pointing to the Game Service"  env:"GAME_URL" required:"true"`
	PodName    string `help:"K8s pod name" env:"POD_NAME" default:"unknown"`
	Namespace  string `help:"K8s namespace" env:"NAMESPACE" default:"unknown"`
	LogLevel   string `help:"Log level (debug, info, warn, error)"    env:"LOG_LEVEL" default:"info"`
}

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

func (s *Config) GetJWTPublic() string {
	return s.JWTPublic
}

func (s *Config) GetPort() string {
	return s.Port
}

func (s *Config) GetRedisURL() string {
	return s.RedisURL
}

func (s *Config) GetPlayersURL() string {
	return s.PlayersURL
}

func (s *Config) GetGameURL() string {
	return s.GameURL
}

func (s *Config) GetPodName() string {
	return s.PodName
}

func (s *Config) GetNamespace() string {
	return s.Namespace
}

func (s *Config) GetLogLevel() string {
	return s.LogLevel
}
