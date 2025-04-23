package handlers

import (
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/MommusWinner/MicroDurak/services/matchmaker/types"
	"github.com/labstack/echo/v4"
)

func AddRoutes(
	e *echo.Echo,
	queue chan<- types.MatchChan,
	config *config.Config,
	playersClient players.PlayersClient,
) {
	h := Handler{
		Queue:         queue,
		Config:        config,
		PlayersClient: playersClient,
	}
	e.GET("/find-match", h.FindMatch)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
