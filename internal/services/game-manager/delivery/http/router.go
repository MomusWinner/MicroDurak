package http

import (
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/lib/jwt"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, handler *GameManagerHandler, ctx domain.Context) {
	e.GET("/api/v1/game-manager/:gameId", handler.Connect, jwt.AuthMiddleware(ctx.Config().GetJWTPublic()))
}
