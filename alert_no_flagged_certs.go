package main

import (
	"fmt"
)

func alertNoFlaggedCerts(secretCount, certCount int) error {
	message := SlackMessage{}

	secretText := "secret"
	if secretCount > 1 {
		secretText += "s"
	}

	certText := "cert"
	if certCount > 1 {
		certText += "s"
	}

	title := fmt.Sprintf(
		"No flagged certificates found!\n Scanned *%v %s* and found *%v %s*",
		secretCount,
		secretText,
		certCount,
		certText,
	)

	message.Blocks = append(message.Blocks, SlackSection{
		Type: "section",
		Text: SlackBlock{
			Type: "mrkdwn",
			Text: title,
		},
	})

	return sendToSlack(message)
}
