package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func AddRoutes(e *echo.Echo, playerHandler *PlayerHandler) {
	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	e.GET("/api/v1/players", playerHandler.GetAll)
	e.GET("/api/v1/players/:id", playerHandler.GetById)
}
