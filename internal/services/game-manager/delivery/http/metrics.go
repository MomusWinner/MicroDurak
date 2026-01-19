package http

import (
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/metrics"
)

type metricsAdapter struct{}

func NewMetricsAdapter() domain.Metrics {
	return &metricsAdapter{}
}

func (m *metricsAdapter) IncPlayersConnected(podName, namespace string) {
	metrics.PlayersConnected.WithLabelValues(podName, namespace).Inc()
}

func (m *metricsAdapter) DecPlayersConnected(podName, namespace string) {
	metrics.PlayersConnected.WithLabelValues(podName, namespace).Dec()
}
