package config

import (
	"github.com/alecthomas/kong"
	"log"

	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/infra"
)

type config struct {
	JWTPublic   string `help:"Base64 Private key for the jwt"       env:"JWT_PUBLIC"   required:"true"`
	RabbitmqURL string `help:"Rabbitmq connection URL"              env:"RABBITMQ_URL" required:"true"`
	Port        string `help:"Port to listen on"                    env:"PORT"                         default:"7070"`
	PodName     string `help:"K8s pod name" env:"POD_NAME" default:"unknown"`
	Namespace   string `help:"K8s namespace" env:"NAMESPACE" default:"unknown"`
	LogLevel    string `help:"Log level (debug, info, warn, error)" env:"LOG_LEVEL"                    default:"info"`
}

func Make() infra.Config {
	cfg := &config{}
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

func (s *config) GetJWTPublic() string {
	return s.JWTPublic
}

func (s *config) GetRabbitmqURL() string {
	return s.RabbitmqURL
}

func (s *config) GetPort() string {
	return s.Port
}

func (s *config) GetPodName() string {
	return s.PodName
}

func (s *config) GetNamespace() string {
	return s.Namespace
}

func (s *config) GetLogLevel() string {
	return s.LogLevel
}
