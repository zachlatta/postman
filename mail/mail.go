// Adapted from the Google App Engine github.com/scorredoira/email packages.
package mail

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"text/template"

	"github.com/jpoehls/gophermail"
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
	msg *gophermail.Message
}

func NewMessage(from, to *mail.Address, subject, templatePath,
	htmlTemplatePath string, context interface{}) (*message, error) {
	msg := &message{
		msg: &gophermail.Message{
			From:    *from,
			To:      []mail.Address{*to},
			Subject: subject,
		},
	}

	if templatePath != "" {
		parsed, err := parseTemplate(templatePath, context)
		if err != nil {
			return nil, err
		}

		msg.msg.Body = parsed
	}

	if htmlTemplatePath != "" {
		parsed, err := parseTemplate(htmlTemplatePath, context)
		if err != nil {
			return nil, err
		}

		msg.msg.HTMLBody = parsed
	}

	return msg, nil
}

func parseTemplate(templatePath string, context interface{}) (string, error) {
	tmplBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("emailBody").Parse(string(tmplBytes)))

	var doc bytes.Buffer
	err = t.Execute(&doc, context)
	if err != nil {
		return "", err
	}

	return string(doc.Bytes()), nil
}

// Send sends an email message.
func (m *Mailer) Send(msg *message) error {
	err := gophermail.SendMail(
		m.Address,
		m.Auth,
		msg.msg,
	)
	if err != nil {
		return errors.New("Error sending email: " + err.Error())
	}

	return nil
}
