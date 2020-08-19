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

type SlackSection struct {
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

	message.Blocks = append(message.Blocks, SlackSection{
		Type: "header",
		Text: SlackBlock{
			Type: "plain_text",
			Text: fmt.Sprintf("You have certificates expiring within %v days", expiresInDays),
		},
	})

	// Add our divider block
	message.Blocks = append(message.Blocks, SlackBlock{
		Type: "divider",
	})

	// group our certs by their fingerprint
	groupedCerts := groupCerts(certs)

	// Add each cert to our message
	for _, group := range groupedCerts {
		if len(group.Certs) > 0 {
			message.Blocks = append(message.Blocks, group.ToSlackMessage())
		}
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
