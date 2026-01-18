package core

import (
	"embed"
	"html/template"

	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/infra"
	"github.com/MommusWinner/MicroDurak/lib/smtp"
)

//go:embed templates/*.html
var templatesFS embed.FS

type SmtpGreeting struct {
	smtp      smtp.SMTP
	templates *template.Template
}

func MakeSMTP(cfg infra.Config) *SmtpGreeting {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic(err)
	}

	smtp := smtp.SMTP{
		EmailFrom: cfg.GetEmailFrom(),
		User:      cfg.GetSMTPUser(),
		Pass:      cfg.GetSMTPPass(),
		Host:      cfg.GetSMTPHost(),
		Port:      cfg.GetSMTPPort(),
	}

	return &SmtpGreeting{
		smtp:      smtp,
		templates: tmpl,
	}
}

func (s *SmtpGreeting) Send(email string, name string) error {
	return s.smtp.SendEmailWithTemplate(email, &smtp.EmailData{
		Name:    name,
		Subject: "Your account verification code",
	}, s.templates, "greeting.html")
}
