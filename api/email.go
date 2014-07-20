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
	sender             MailSender
	config             *AppConfig
	defaultFromAddress string
	templates          map[string]*template.Template
}

func (m *Mailer) Send(msg *Message) error {
	return m.sender.Send(msg)
}

func (m *Mailer) ParseTemplate(name string) (*template.Template, error) {
	var (
		t   *template.Template
		ok  bool
		err error
	)
	t, ok = m.templates[name]
	if !ok {
		t, err = template.ParseFiles(path.Join(m.config.TemplatesDir, name+".tmpl"))
		if err != nil {
			return nil, err
		}
		m.templates[name] = t

	}
	return t, nil
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
	t, err := m.ParseTemplate(templateName)
	if err != nil {
		return nil, err
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

func (s *smtpSender) Send(msg *Message) error {
	return smtp.SendMail(s.config.SmtpHost+":25", s.Auth, msg.From, msg.To, msg.Body)
}

type fakeSender struct{}

func (s *fakeSender) Send(msg *Message) error {
	log.Println(msg)
	return nil
}

func newSmtpSender(config *AppConfig) *smtpSender {
	s := &smtpSender{config: config}
	s.Auth = smtp.PlainAuth("",
		config.SmtpName,
		config.SmtpPassword,
		config.SmtpHost,
	)
	return s
}

func NewMailer(config *AppConfig) *Mailer {
	mailer := &Mailer{config: config}
	if config.SmtpName == "" {
		log.Println("WARNING: using fake mailer, messages will not be sent by SMTP. " +
			"Set SMTP_NAME and SMTP_PASSWORD in environment to enable.")
		mailer.sender = &fakeSender{}
	} else {
		mailer.sender = newSmtpSender(config)
	}
	mailer.defaultFromAddress = config.SmtpDefaultSender
	mailer.templates = make(map[string]*template.Template)
	return mailer
}

func (m *Mailer) SendResetPasswordMail(user *User, recoveryCode string, r *http.Request) error {
	msg, err := m.MessageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		m.defaultFromAddress,
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

func (m *Mailer) SendWelcomeMail(user *User) error {
	msg, err := m.MessageFromTemplate(
		"Welcome to photoshare!",
		[]string{user.Email},
		m.defaultFromAddress,
		"signup",
		user,
	)
	if err != nil {
		return err
	}
	return m.Send(msg)
}
