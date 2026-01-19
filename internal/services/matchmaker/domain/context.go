package domain

import (
	"log/slog"

	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/infra"
)

type Context interface {
	Make() Context
	Connection() Connection
	Config() infra.Config
	Logger() *slog.Logger
}
