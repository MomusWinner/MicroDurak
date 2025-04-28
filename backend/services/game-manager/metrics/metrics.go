package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PlayersConnected = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "game_manager_players_connected",
			Help: "Current players connected to the game",
		},
		[]string{"pod", "namespace"},
	)
)
