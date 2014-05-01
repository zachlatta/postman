package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/zachlatta/postman/mail"
)

type Recipient map[string]string

var (
	htmlTemplatePath, textTemplatePath        string
	csvPath                                   string
	smtpURL, smtpUser, smtpPassword, smtpPort string
	sender, subject                           string
)

var flags []*flag.Flag

func main() {
	flag.StringVar(&htmlTemplatePath, "html", "", "html template path")
	flag.StringVar(&textTemplatePath, "text", "", "text template path")
	flag.StringVar(&csvPath, "csv", "", "path to csv of contact list")
	flag.StringVar(&smtpURL, "server", "", "url of smtp server")
	flag.StringVar(&smtpPort, "port", "", "port of smtp server")
	flag.StringVar(&smtpUser, "user", "", "smtp username")
	flag.StringVar(&smtpPassword, "password", "", "smtp password")
	flag.StringVar(&sender, "sender", "", "email to send from")
	flag.StringVar(&subject, "subject", "", "subject of email")

	flag.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})

	flag.Usage = usage

	flag.Parse()
	log.SetFlags(0)

	checkAndHandleMissingFlags(flags)

	csv, err := os.Open(csvPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening CSV:", err.Error())
		os.Exit(2)
	}
	defer csv.Close()

	recipients, emailField, err := readCSV(csvPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading CSV:", err.Error())
		os.Exit(2)
	}

	mailer := mail.NewMailer(
		smtpUser,
		smtpPassword,
		smtpURL,
		smtpPort,
	)

	success := make(chan Recipient)
	fail := make(chan error)

	go func() {
		for _, recipient := range *recipients {
			go sendMail(recipient, *emailField, &mailer, success, fail)
		}
	}()

	for i := 0; i < len(*recipients); i++ {
		select {
		case <-success:
			fmt.Printf("\rEmailed recipient %d of %d...", i+1, len(*recipients))
		case err := <-fail:
			fmt.Fprintln(os.Stderr, "\nError sending email:", err.Error())
			os.Exit(2)
		}
	}
	fmt.Println()
}

func checkAndHandleMissingFlags(flags []*flag.Flag) {
	var flagsMissing []*flag.Flag
	for _, f := range flags {
		if f.Value.String() == "" {
			flagsMissing = append(flagsMissing, f)
		}
	}

	missingCount := len(flagsMissing)
	if missingCount > 0 {
		if missingCount == len(flags) {
			usage()
		}

		missingFlags(flagsMissing)
	}
}

const usageTemplate = `Postman is a utility for sending batch emails.

Usage:

  postman [flags]

Flags:
{{range .}}
  -{{.Name | printf "%-11s"}} {{.Usage}}{{end}}

`

const missingFlagsTemplate = `Missing required flags:
{{range .}}
  -{{.Name | printf "%-11s"}} {{.Usage}}{{end}}

`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, flags)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printMissingFlags(w io.Writer, missingFlags []*flag.Flag) {
	tmpl(w, missingFlagsTemplate, missingFlags)
}

func missingFlags(missingFlags []*flag.Flag) {
	printMissingFlags(os.Stderr, missingFlags)
	os.Exit(2)
}

func sendMail(recipient Recipient, emailField string, mailer *mail.Mailer,
	success chan Recipient, fail chan error) {

	message := &mail.Message{
		Sender:   sender,
		To:       []string{recipient[emailField]},
		Subject:  subject,
		Template: textTemplatePath,
		Context:  recipient,
	}

	if err := mailer.Send(message); err != nil {
		fail <- err
		return
	}

	success <- recipient
}
