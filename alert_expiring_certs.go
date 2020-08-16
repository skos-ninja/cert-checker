package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SlackMessage struct {
	Blocks []interface{} `json:"blocks"`
}

type SlackHeader struct {
	Type string     `json:"type"`
	Text SlackBlock `json:"text"`
}

type SlackBlock struct {
	Type     string       `json:"type"`
	Text     string       `json:"text,omitempty"`
	Elements []SlackBlock `json:"elements,omitempty"`
}

func alertExpiringCerts(certs []Cert) error {
	if slackWebHook == "" {
		log.Println("No slack webhook set. Skipping alert")
		return nil
	}

	message := SlackMessage{}

	message.Blocks = append(message.Blocks, SlackHeader{
		Type: "header",
		Text: SlackBlock{
			Type: "plain_text",
			Text: fmt.Sprintf("You have certificates expiring within %v days:", expiresInDays),
		},
	})

	// Add our divider block
	message.Blocks = append(message.Blocks, SlackBlock{
		Type: "divider",
	})

	// Add each cert as a block
	for _, cert := range certs {
		message.Blocks = append(message.Blocks, cert.ToSlackMessage())
	}

	// Send our webhook to slack
	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = http.Post(slackWebHook, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
