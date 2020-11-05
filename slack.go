package main

import (
	"bytes"
	"encoding/json"
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

func sendToSlack(message SlackMessage) error {
	if slackWebHook == "" {
		log.Println("No slack webhook set. Skipping alert")
		return nil
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
