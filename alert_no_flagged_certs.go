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

	environmentText := ""
	if environmentString != "" {
		environmentText = fmt.Sprintf("in %s", environmentString)
	}

	title := fmt.Sprintf(
		"No flagged certificates found!\n Scanned *%v %s* and found *%v %s* %s",
		secretCount,
		secretText,
		certCount,
		certText,
		environmentText,
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
