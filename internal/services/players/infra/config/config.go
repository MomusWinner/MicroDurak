package config

import (
	"github.com/alecthomas/kong"
	"log"
)

import ()

type Config struct {
	JWTPublic   string `help:"Base64 Private key for the jwt" env:"JWT_PUBLIC" required:"true"`
	HTTPPort    string `help:"Port to listen on" env:"HTTP_PORT" default:"8080"`
	GRPCPort    string `help:"Port to listen on" env:"GRPC_PORT" default:"9090"`
	DatabaseURL string `help:"Database connection URL" env:"DATABASE_URL" required:"true"`
	LogLevel    string `help:"Log level (debug, info, warn, error)" env:"LOG_LEVEL" default:"info"`
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

func (s *Config) GetJwtPublic() string {
	return s.JWTPublic
}

func (s *Config) GetHTTPPort() string {
	return s.HTTPPort
}

func (s *Config) GetGRPCPort() string {
	return s.GRPCPort
}
func (s *Config) GetDatabaseURL() string {
	return s.DatabaseURL
}

func (s *Config) GetLogLevel() string {
	return s.LogLevel
}
