package domain

import (
	"log/slog"

	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain/infra"
)

type Context interface {
	Make() Context
	Config() infra.Config
	Logger() *slog.Logger
	Messaging() Messaging
}
