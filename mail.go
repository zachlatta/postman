package main

import (
	stdMail "net/mail"

	"github.com/zachlatta/postman/mail"
)

func sendMail(recipient Recipient, emailField string, mailer *mail.Mailer,
	debug bool, success chan *mail.Message, fail chan error) {

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

	message, err := mail.NewMessage(
		parsedSender,
		parsedTo,
		subject,
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
