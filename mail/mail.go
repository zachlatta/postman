// Adapted from Google Appengine mail package.
package mail

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/smtp"
	"text/template"
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

// An email message. Include either the Body or Template field. If both are
// included, only Template will be used.
type Message struct {
	Sender   string
	To       []string
	Subject  string
	Body     string
	Template string      // Path of template to load
	Context  interface{} // required if Template is used
}

const emailTemplate = `From: {{.Sender}}
To: {{range .To}}{{.}},{{end}}
Subject: {{.Subject}}

{{.Body}}
`

// Send sends an email message.
func (m *Mailer) Send(msg *Message) error {
	if msg.Template != "" {
		tmplBytes, err := ioutil.ReadFile(msg.Template)
		if err != nil {
			return err
		}

		t := template.Must(template.New("emailBody").Parse(string(tmplBytes)))

		var doc bytes.Buffer
		err = t.Execute(&doc, msg.Context)
		if err != nil {
			return err
		}

		msg.Body = string(doc.Bytes())
	}

	return m.send(msg)
}

func (m *Mailer) send(msg *Message) error {
	var doc bytes.Buffer

	t := template.New("emailTemplate")
	t, err := t.Parse(emailTemplate)
	if err != nil {
		return errors.New("Error parsing mail template: " + err.Error())
	}
	err = t.Execute(&doc, msg)
	if err != nil {
		return errors.New("Error executing mail template: " + err.Error())
	}

	err = smtp.SendMail(
		m.Address,
		m.Auth,
		msg.Sender,
		msg.To,
		doc.Bytes(),
	)
	if err != nil {
		return errors.New("Error sending email: " + err.Error())
	}

	return nil
}
