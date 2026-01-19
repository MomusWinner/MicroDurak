package infra

type Config interface {
	GetJWTPublic() string
	GetPort() string
	GetRedisURL() string
	GetPlayersURL() string
	GetGameURL() string
	GetPodName() string
	GetNamespace() string
	GetLogLevel() string
}
