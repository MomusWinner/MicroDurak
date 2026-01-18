package smtp

import (
	"bytes"
	"crypto/tls"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"html/template"
)

type EmailData struct {
	Name    string
	Subject string
}

type SMTP struct {
	EmailFrom   string
	User        string
	Pass        string
	Host        string
	Port        int
	TemplateDir string
}

func (s *SMTP) SendEmailWithTemplate(email string, data *EmailData, tmpl *template.Template, emailTemp string) error {
	from := s.EmailFrom
	smtpPass := s.Pass
	smtpUser := s.User
	to := email
	smtpHost := s.Host
	smtpPort := s.Port

	var body bytes.Buffer

	if err := tmpl.ExecuteTemplate(&body, emailTemp, data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
