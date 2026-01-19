package infra

type Config interface {
	GetJWTPublic() string
	GetRabbitmqURL() string
	GetPort() string
	GetPodName() string
	GetNamespace() string
	GetLogLevel() string
}
