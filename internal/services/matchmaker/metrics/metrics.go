package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PlayersSearching = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "matchmaker_players_searching",
			Help: "Current players in the matchmaker",
		},
		[]string{"pod", "namespace"},
	)

	SearchDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "matchmaker_search_duration_avg",
			Help: "Duration of matchmaking search",
		},
		[]string{"pod", "namespace"},
	)

	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "matchmaker_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status_code", "pod", "namespace"},
	)

	WebsocketUpgradeErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "matchmaker_websocket_upgrade_errors_total",
			Help: "Total number of WebSocket upgrade failures",
		},
		[]string{"pod", "namespace"},
	)

	WebsocketWriteErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "matchmaker_websocket_write_errors_total",
			Help: "Total number of WebSocket write failures",
		},
		[]string{"pod", "namespace"},
	)
)
