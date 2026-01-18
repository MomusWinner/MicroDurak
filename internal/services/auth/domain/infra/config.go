package infra

type Config interface {
	GetJwtPrivate() string
	GetPlayersURL() string
	GetPort() string
	GetDatabaseURL() string
	GetLogLevel() string
	GetEmailFrom() string
	GetSMTPHost() string
	GetSMTPPass() string
	GetSMTPPort() int
	GetSMTPUser() string
}
