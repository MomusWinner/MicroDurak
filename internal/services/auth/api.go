package auth

import (
	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/config"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/handlers"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, config *config.Config, queries *database.Queries, playersClient players.PlayersClient) {
	h := handlers.Handler{DBQueries: queries, Config: config, PlayersClient: playersClient}
	e.POST("/auth/register", h.Register)
	e.POST("/auth/login", h.Login)
}
