package main

import (
	"fmt"
	"time"
)

func alertExpiringCerts(certs []Cert, t time.Time) error {
	message := SlackMessage{}

	// group our certs by their fingerprint
	groupedCerts := groupCerts(certs, t)

	title := ""
	if groupedCerts[0].ExpiresInDays < 1 {
		title = fmt.Sprintf("You have EXPIRED CERTIFICATES")
		if environmentString != "" {
			title += fmt.Sprintf(" in %s", environmentString)
		}
		title += "!"
	} else {
		days := "days"
		if groupedCerts[0].ExpiresInDays == 1 {
			days = "day"
		}
		title = fmt.Sprintf("You have certificates expiring within %v %s", groupedCerts[0].ExpiresInDays, days)
		if environmentString != "" {
			title += fmt.Sprintf(" in %s", environmentString)
		}
	}
	message.Blocks = append(message.Blocks, SlackSection{
		Type: "header",
		Text: SlackBlock{
			Type: "plain_text",
			Text: title,
		},
	})

	// Add our divider block
	message.Blocks = append(message.Blocks, SlackBlock{
		Type: "divider",
	})

	// Add each cert to our message
	for _, group := range groupedCerts {
		if len(group.Certs) > 0 {
			message.Blocks = append(message.Blocks, group.ToSlackMessage())
		}
	}

	return sendToSlack(message)
}
