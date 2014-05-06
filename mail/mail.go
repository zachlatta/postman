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

func NewMessage(from, to *mail.Address, subject string, files []string, templatePath,
	htmlTemplatePath string, context interface{}) (*email.Email, error) {
	msg := &email.Email{
		From:    from.String(),
		To:      []string{to.String()},
		Subject: subject,
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

	t := template.Must(template.New("emailBody").Parse(string(tmplBytes)))

	var doc bytes.Buffer
	err = t.Execute(&doc, context)
	if err != nil {
		return nil, err
	}

	return doc.Bytes(), nil
}

// Send sends an email Message.
func (m *Mailer) Send(msg *email.Email) error {
	err := msg.Send(
		m.Address,
		m.Auth,
	)
	if err != nil {
		return errors.New("Error sending email: " + err.Error())
	}

	return nil
}
