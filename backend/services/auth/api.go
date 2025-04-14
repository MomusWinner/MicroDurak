package auth

import (
	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/MommusWinner/MicroDurak/services/auth/handlers"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, config *config.Config, queries *database.Queries) {
	h := handlers.Handler{DBQueries: queries, Config: config}
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
}
