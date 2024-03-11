package main

import (
	stdMail "net/mail"

	"github.com/pepa65/postman/mail"
	"gopkg.in/jordan-wright/email.v2"
)

func sendMail(recipient Recipient, emailField string, mailer *mail.Mailer,
	debug bool, success chan *email.Email, fail chan error) {

	parsedSender, err := stdMail.ParseAddress(sender)
	if err != nil {
		fail <- err
		return
	}

	parsedTo, err := stdMail.ParseAddress(recipient[emailField])
	if err != nil {
		fail <- err
		return
	}

	// Should probably add a clearer distinction somewhere, but in the case that
	// the user doesn't provide the -html flag, htmlTemplatePath will be an empty
	// string. If it's an empty string, then it'll be ignored and not parsed and
	// added to the message within the mail.NewMessage method.
	message, err := mail.NewMessage(
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
