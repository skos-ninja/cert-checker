package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	resp, err := http.Post(slackWebHook, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		log.Printf("Error sending message to slack  %s", err.Error())
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		return fmt.Errorf("error sending message to slack. status code: %d. resp: %s", resp.StatusCode, bodyString)
	}

	return nil
}
