package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func AddRoutes(e *echo.Echo, playerHandler *PlayerHandler) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/api/v1/players", playerHandler.GetAll)
	e.GET("/api/v1/players/:id", playerHandler.GetById)
	e.POST("/api/v1/matches", playerHandler.CreateMatch)
	e.GET("/api/v1/matches/:id", playerHandler.GetMatchResultById)
	e.GET("/api/v1/matches", playerHandler.GetAllMatchResults)
}
