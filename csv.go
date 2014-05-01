package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

var (
	missingEmailField = errors.New("Email field missing in header.")
)

func readCSV(path string) (*[]Recipient, *string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var (
		header     []string
		headerRead bool
		emailField string
		recipients []Recipient
	)

	reader := csv.NewReader(file)
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
				return nil, nil, missingEmailField
			}

			headerRead = true
		}
	}

	return &recipients, &emailField, nil
}
