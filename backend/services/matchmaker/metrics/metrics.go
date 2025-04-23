package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PlayersSearching = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "matchmaker_players_searching",
		Help: "Current players in the matchmaker",
	})

	SearchDuration = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "matchmaker_search_duration_avg",
		Help: "Duration of matchmaking search",
	})

	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "matchmaker_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	WebsocketUpgradeErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "matchmaker_websocket_upgrade_errors_total",
			Help: "Total number of WebSocket upgrade failures",
		},
	)

	WebsocketWriteErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "matchmaker_websocket_write_errors_total",
			Help: "Total number of WebSocket write failures",
		},
	)
)
