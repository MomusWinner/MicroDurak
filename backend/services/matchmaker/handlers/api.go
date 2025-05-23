package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/lib/jwt"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/config"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/types"
)

func AddRoutes(
	e *echo.Echo,
	queue chan<- types.MatchChan,
	cancel chan<- types.MatchCancel,
	config *config.Config,
	playersClient players.PlayersClient,
) {
	h := Handler{
		Queue:         queue,
		Cancel:        cancel,
		Config:        config,
		PlayersClient: playersClient,
	}
	e.GET("/matchmaker/find-match", h.FindMatch, jwt.AuthMiddleware(config.JWTPublic))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
