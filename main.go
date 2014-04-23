package main

import (
	"flag"
	"io"
	"log"
	"os"
	"pkg/text/template"
)

var (
	htmlTemplatePath, textTemplatePath        string
	csvPath                                   string
	smtpURL, smtpUser, smtpPassword, smtpPort string
)

var flags []*flag.Flag

func main() {
	flag.StringVar(&htmlTemplatePath, "html", "", "html template path")
	flag.StringVar(&textTemplatePath, "text", "", "text template path")
	flag.StringVar(&csvPath, "csv", "", "path to csv of contact list")
	flag.StringVar(&smtpURL, "server", "", "url of smtp server")
	flag.StringVar(&smtpUser, "user", "", "smtp username")
	flag.StringVar(&smtpPassword, "password", "", "smtp password")

	flag.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})

	flag.Usage = usage

	flag.Parse()
	log.SetFlags(0)

	checkAndHandleMissingFlags(flags)
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
