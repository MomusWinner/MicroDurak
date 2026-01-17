package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func AddRoutes(e *echo.Echo, authHandler *AuthHandler) {
	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	e.POST("/api/v1/auth/login", authHandler.Login)
	e.POST("/api/v1/auth/register", authHandler.Register)
}
