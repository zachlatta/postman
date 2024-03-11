// Adapted from the Google App Engine github.com/scorredoira/email packages.
package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"text/template"

	"gopkg.in/jordan-wright/email.v2"
)

// Mailer encapsulates data used for sending email.
type Mailer struct {
	Auth    smtp.Auth
	Address string
	TLS     *tls.Config
}

// NewMailer creates a new Mailer.
func NewMailer(username, password, host, port string, skipCertValidation bool) Mailer {
	return Mailer{
		Auth: smtp.PlainAuth(
			"",
			username,
			password,
			host,
		),
		Address: host + ":" + port,
		TLS: &tls.Config{
			InsecureSkipVerify: skipCertValidation,
			ServerName:         host,
		},
	}
}

func NewMessage(from, to *mail.Address, subject string, files []string, templatePath,
	htmlTemplatePath string, context interface{}) (*email.Email, error) {
	parsedSubject, err := parseText([]byte(subject), context)
	if err != nil {
		return nil, err
	}

	msg := &email.Email{
		From:    from.String(),
		To:      []string{to.String()},
		Subject: string(parsedSubject),
	}

	for _, file := range files {
		_, err := msg.AttachFile(file)
		if err != nil {
			return nil, err
		}
	}

	if templatePath != "" {
		parsed, err := parseTemplate(templatePath, context)
		if err != nil {
			return nil, err
		}

		msg.Text = parsed
	}

	if htmlTemplatePath != "" {
		parsed, err := parseTemplate(htmlTemplatePath, context)
		if err != nil {
			return nil, err
		}

		msg.HTML = parsed
	}

	return msg, nil
}

func parseTemplate(templatePath string, context interface{}) ([]byte, error) {
	tmplBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	return parseText(tmplBytes, context)
}

func parseText(tmplBytes []byte, context interface{}) ([]byte, error) {
	t := template.Must(template.New("emailBody").Parse(string(tmplBytes)))
	var txt bytes.Buffer
	err := t.Execute(&txt, context)
	if err != nil {
		return nil, err
	}

	return txt.Bytes(), nil
}

// Send sends an email Message.
func (m *Mailer) Send(msg *email.Email) error {
	err := msg.SendWithTLS(
		m.Address,
		m.Auth,
		m.TLS,
	)
	if err != nil {
		return errors.New("Error sending email: " + err.Error())
	}

	return nil
}
