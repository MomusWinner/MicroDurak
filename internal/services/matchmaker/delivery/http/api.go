package http

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/config"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/types"
	"github.com/MommusWinner/MicroDurak/lib/jwt"
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
	e.GET("/api/v1/matchmaker/find-match", h.FindMatch, jwt.AuthMiddleware(config.JWTPublic))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
