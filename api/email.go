package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"path"
	"strings"
	"text/template"
)

type Message struct {
	Subject string
	Body    []byte
	To      []string
	From    string
}

func (msg *Message) String() string {
	return fmt.Sprintf("To:%s\r\nFrom:%s\r\nSubject: %s\r\nBody:%s",
		strings.Join(msg.To, ", "),
		msg.From,
		msg.Subject,
		string(msg.Body),
	)
}

type Mailer struct {
	Sender    MailSender
	Config    *AppConfig
	Templates map[string]*template.Template
}

func (m *Mailer) Send(msg *Message) error {
	return m.Sender.Send(msg)
}

func (m *Mailer) ParseTemplate(name string) (*template.Template, error) {
	return template.ParseFiles(path.Join(m.Config.TemplatesDir, name))
}

// Creates a new message from a template; message body set to rendered template
func (m *Mailer) MessageFromTemplate(subject string,
	to []string,
	from string,
	templateName string,
	data interface{}) (*Message, error) {

	msg := &Message{
		Subject: subject,
		To:      to,
		From:    from,
	}
	b := &bytes.Buffer{}
	t, ok := m.Templates[templateName]
	if !ok {
		t, err := m.ParseTemplate(templateName + ".tmpl")
		if err != nil {
			return nil, err
		}
		m.Templates[templateName] = t
	}
	if err := t.Execute(b, data); err != nil {
		return nil, err
	}
	msg.Body = b.Bytes()
	return msg, nil
}

type MailSender interface {
	Send(*Message) error
}

type smtpSender struct {
	smtp.Auth
	config *AppConfig
}

func (m *smtpSender) Send(msg *Message) error {
	return smtp.SendMail(m.config.SmtpHost+":25", m.Auth, msg.From, msg.To, msg.Body)
}

type fakeSender struct{}

func (m *fakeSender) Send(msg *Message) error {
	log.Println(msg)
	return nil
}

func newSmtpSender(config *AppConfig) *smtpSender {
	s := &smtpSender{config: config}
	s.Auth = smtp.PlainAuth("", config.SmtpName, config.SmtpPassword, config.SmtpHost)
	return s
}

func NewMailer(config *AppConfig) *Mailer {
	mailer := &Mailer{Config: config}
	if config.SmtpName == "" {
		log.Println("WARNING: using fake mailer, messages will not be sent by SMTP. " +
			"Set SMTP_NAME and SMTP_PASSWORD in environment to enable.")
		mailer.Sender = &fakeSender{}
	} else {
		mailer.Sender = newSmtpSender(config)
	}
	mailer.Templates = make(map[string]*template.Template)
	return mailer
}

func (m *Mailer) sendResetPasswordMail(user *User, recoveryCode string, r *http.Request) error {
	msg, err := m.MessageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		m.Config.SmtpDefaultSender,
		"recover_pass",
		&struct {
			Name         string
			RecoveryCode string
			Url          string
		}{
			user.Name,
			recoveryCode,
			baseURL(r),
		},
	)
	if err != nil {
		return err
	}
	return m.Send(msg)
}

func (m *Mailer) sendWelcomeMail(user *User) error {
	msg, err := m.MessageFromTemplate(
		"Welcome to photoshare!",
		[]string{user.Email},
		m.Config.SmtpDefaultSender,
		"signup",
		user,
	)
	if err != nil {
		return err
	}
	return m.Send(msg)
}
