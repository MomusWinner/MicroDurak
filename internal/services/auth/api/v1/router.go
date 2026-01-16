package v1

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, authHandler *AuthHandler) {
	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/register", authHandler.Register)
}
