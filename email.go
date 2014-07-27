package photoshare

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

type message struct {
	subject string
	body    []byte
	to      []string
	from    string
}

func (msg *message) String() string {
	return fmt.Sprintf("From:%s\r\nTo:%s\r\nSubject: %s\r\n\r\n%s",
		msg.from,
		strings.Join(msg.to, ", "),
		msg.subject,
		string(msg.body),
	)
}

type mailer struct {
	sender             mailSender
	config             *appConfig
	defaultFromAddress string
	templates          map[string]*template.Template
}

func (m *mailer) send(msg *message) error {
	return m.sender.send(msg)
}

func (m *mailer) parseTemplate(name string) (*template.Template, error) {
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
func (m *mailer) messageFromTemplate(subject string,
	to []string,
	from string,
	templateName string,
	data interface{}) (*message, error) {

	msg := &message{
		subject: subject,
		to:      to,
		from:    from,
	}
	b := &bytes.Buffer{}
	t, err := m.parseTemplate(templateName)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(b, data); err != nil {
		return nil, errgo.Mask(err)
	}
	msg.body = b.Bytes()
	return msg, nil
}

type mailSender interface {
	send(*message) error
}

type smtpSender struct {
	smtp.Auth
	config *appConfig
}

func (s *smtpSender) send(msg *message) error {
	return errgo.Mask(smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SmtpHost, s.config.SmtpPort),
		s.Auth,
		msg.from,
		msg.to,
		[]byte(msg.String()),
	))
}

type fakeSender struct{}

func (s *fakeSender) send(msg *message) error {
	log.Println(msg)
	return nil
}

func newSmtpSender(config *appConfig) *smtpSender {
	s := &smtpSender{config: config}
	s.Auth = smtp.PlainAuth("",
		config.SmtpName,
		config.SmtpPassword,
		config.SmtpHost,
	)
	return s
}

func newMailer(config *appConfig) *mailer {
	mailer := &mailer{config: config}
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

func (m *mailer) sendResetPasswordMail(user *user, recoveryCode string, r *http.Request) error {
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
			getBaseURL(r),
		},
	)
	if err != nil {
		return err
	}
	return m.send(msg)
}

func (m *mailer) sendWelcomeMail(user *user) error {
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
	return m.send(msg)
}
