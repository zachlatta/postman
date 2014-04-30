package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
	defer csv.Close()

	recipients, emailField, err := readCSV(csvPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(2)
	}

	mailer := mail.NewMailer(
		smtpUser,
		smtpPassword,
		smtpURL,
		smtpPort,
	)

	// TODO: Make concurrent
	for _, recipient := range *recipients {
		message := &mail.Message{
			Sender:   sender,
			To:       []string{recipient[*emailField]},
			Subject:  subject,
			Template: textTemplatePath,
			Context:  recipient,
		}

		if err := mailer.Send(message); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(2)
		}
	}
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

func readCSV(path string) (*[]Recipient, *string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var (
		header     []string
		emailField string
		recipients []Recipient
	)

	reader, headerRead := csv.NewReader(file), false
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}

		if headerRead {
			recipient := make(Recipient)

			for i, key := range header {
				recipient[key] = fields[i]
			}

			recipients = append(recipients, recipient)
		} else {
			header = fields

			for _, v := range header {
				if strings.ToLower(v) == "email" {
					emailField = v
				}
			}

			if emailField == "" {
				return nil, nil, errors.New("Email field missing in header.")
			}

			headerRead = true
		}
	}

	return &recipients, &emailField, nil
}
