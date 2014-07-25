package api

import (
	"bytes"
	"fmt"
	"github.com/juju/errgo"
	"log"
	"net/http"
	"net/smtp"
	"path"
	"strings"
	"text/template"
)

// Message models an email message
type Message struct {
	Subject string
	Body    []byte
	To      []string
	From    string
}

func (msg *Message) String() string {
	return fmt.Sprintf("From:%s\r\nTo:%s\r\nSubject: %s\r\n\r\n%s",
		msg.From,
		strings.Join(msg.To, ", "),
		msg.Subject,
		string(msg.Body),
	)
}

// Mailer keeps email sender and template info
type Mailer struct {
	sender             MailSender
	config             *AppConfig
	defaultFromAddress string
	templates          map[string]*template.Template
}

// Send sends the message
func (m *Mailer) Send(msg *Message) error {
	return m.sender.Send(msg)
}

func (m *Mailer) parseTemplate(name string) (*template.Template, error) {
	var (
		t   *template.Template
		ok  bool
		err error
	)
	t, ok = m.templates[name]
	if !ok {
		t, err = template.ParseFiles(path.Join(m.config.TemplatesDir, name+".tmpl"))
		if err != nil {
			return nil, errgo.Mask(err)
		}
		m.templates[name] = t

	}
	return t, nil
}

// Creates a new message from a template; message body set to rendered template
func (m *Mailer) messageFromTemplate(subject string,
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
	t, err := m.parseTemplate(templateName)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(b, data); err != nil {
		return nil, errgo.Mask(err)
	}
	msg.Body = b.Bytes()
	return msg, nil
}

// MailSender handles sending of email
type MailSender interface {
	Send(*Message) error
}

type smtpSender struct {
	smtp.Auth
	config *AppConfig
}

func (s *smtpSender) Send(msg *Message) error {
	return errgo.Mask(smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SmtpHost, s.config.SmtpPort),
		s.Auth,
		msg.From,
		msg.To,
		[]byte(msg.String()),
	))
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

// NewMailer creates a new Mailer instance
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

// SendResetPasswordMail sends email to user to reset their password
func (m *Mailer) SendResetPasswordMail(user *User, recoveryCode string, r *http.Request) error {
	msg, err := m.messageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		m.defaultFromAddress,
		"recover_pass",
		&struct {
			Name         string
			RecoveryCode string
			URL          string
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

// SendWelcomeMail sends a welcome message to user
func (m *Mailer) SendWelcomeMail(user *User) error {
	msg, err := m.messageFromTemplate(
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
