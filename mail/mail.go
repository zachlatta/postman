// Adapted from the Google App Engine github.com/scorredoira/email packages.
package mail

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"text/template"

	"github.com/jordan-wright/email"
)

// Mailer encapsulates data used for sending email.
type Mailer struct {
	Auth    smtp.Auth
	Address string
}

// NewMailer creates a new Mailer.
func NewMailer(username, password, host, port string) Mailer {
	return Mailer{
		Auth: smtp.PlainAuth(
			"",
			username,
			password,
			host,
		),
		Address: host + ":" + port,
	}
}

// An email message.
type message struct {
	email *email.Email
}

func NewMessage(from, to *mail.Address, subject, templatePath,
	htmlTemplatePath string, context interface{}) (*message, error) {
	msg := &message{
		email: &email.Email{
			From:    from.String(),
			To:      []string{to.String()},
			Subject: subject,
		},
	}

	if templatePath != "" {
		parsed, err := parseTemplate(templatePath, context)
		if err != nil {
			return nil, err
		}

		msg.email.Text = parsed
	}

	if htmlTemplatePath != "" {
		parsed, err := parseTemplate(htmlTemplatePath, context)
		if err != nil {
			return nil, err
		}

		msg.email.HTML = parsed
	}

	return msg, nil
}

func parseTemplate(templatePath string, context interface{}) ([]byte, error) {
	tmplBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	t := template.Must(template.New("emailBody").Parse(string(tmplBytes)))

	var doc bytes.Buffer
	err = t.Execute(&doc, context)
	if err != nil {
		return nil, err
	}

	return doc.Bytes(), nil
}

// Send sends an email message.
func (m *Mailer) Send(msg *message) error {
	err := msg.email.Send(
		m.Address,
		m.Auth,
	)
	if err != nil {
		return errors.New("Error sending email: " + err.Error())
	}

	return nil
}
