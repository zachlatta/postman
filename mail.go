package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"html/template"
	"io/ioutil"
	"net/mail"
	"net/smtp"

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

func sendMail(recipient Recipient, emailField string, mailer *Mailer,
	debug bool, success chan *email.Email, fail chan error) {

	parsedSender, err := mail.ParseAddress(sender)
	if err != nil {
		fail <- err
		return
	}

	parsedTo, err := mail.ParseAddress(recipient[emailField])
	if err != nil {
		fail <- err
		return
	}

	// Should probably add a clearer distinction somewhere, but in the case that
	// the user doesn't provide the -html flag, htmlTemplatePath will be an empty
	// string. If it's an empty string, then it'll be ignored and not parsed and
	// added to the message within the mail.NewMessage method.
	message, err := NewMessage(
		parsedSender,
		parsedTo,
		subject,
		files,
		textTemplatePath,
		htmlTemplatePath,
		recipient,
	)
	if err != nil {
		fail <- err
		return
	}

	if !debug {
		if err := mailer.Send(message); err != nil {
			fail <- err
			return
		}
	}

	success <- message
}
