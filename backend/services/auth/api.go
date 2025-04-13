package auth

import (
	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/MommusWinner/MicroDurak/services/auth/handlers"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, config *config.Config, queries *database.Queries) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ac := &handlers.AuthContext{
				Context:   c,
				Config:    config,
				DBQueries: queries,
			}
			return next(ac)
		}
	})

	e.POST("/register", handlers.RegisterHandler)
	e.POST("/login", handlers.LoginHandler)
}
