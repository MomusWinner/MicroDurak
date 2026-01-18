package domain

type SMTP interface {
	Send(email string, name string) error
}
