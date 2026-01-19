package http

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/MommusWinner/MicroDurak/lib/jwt"
)

func AddRoutes(
	e *echo.Echo,
	handler *Handler,
) {
	e.GET("/api/v1/matchmaker/find-match", handler.FindMatch, jwt.AuthMiddleware(handler.Ctx.Config().GetJWTPublic()))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
