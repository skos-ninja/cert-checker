package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
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
	groupedCert := make(map[string][]Cert)
	for _, cert := range certs {
		cert := cert
		fingerprint := getMD5Hash(cert.X509.Raw)
		groupedCert[fingerprint] = append(groupedCert[fingerprint], cert)
	}

	// Add each cert to our message
	for fingerprint, certs := range groupedCert {
		if len(certs) > 0 {
			message.Blocks = append(message.Blocks, certsToSlackMessage(fingerprint, certs))
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

func getMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
