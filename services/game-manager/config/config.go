package config

import (
	"github.com/alecthomas/kong"
)

type Config struct {
	JWTPublic   string `help:"Base64 Private key for the jwt"       env:"JWT_PUBLIC"   required:"true"`
	RabbitmqURL string `help:"Rabbitmq connection URL"              env:"RABBITMQ_URL" required:"true"`
	Port        string `help:"Port to listen on"                    env:"PORT"                         default:"7070"`
	PodName     string `help:"K8s pod name" env:"POD_NAME" default:"unknown"`
	Namespace   string `help:"K8s namespace" env:"NAMESPACE" default:"unknown"`
	LogLevel    string `help:"Log level (debug, info, warn, error)" env:"LOG_LEVEL"                    default:"info"`
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
