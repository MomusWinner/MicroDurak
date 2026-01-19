package domain

type Metrics interface {
	IncPlayersConnected(podName, namespace string)
	DecPlayersConnected(podName, namespace string)
}
