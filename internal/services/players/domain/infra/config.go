package infra

type Config interface {
	GetJwtPublic() string
	GetHTTPPort() string
	GetGRPCPort() string
	GetDatabaseURL() string
	GetLogLevel() string
}
