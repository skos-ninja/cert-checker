package main

import (
	"strconv"
	"strings"
)

func certsToSlackMessage(fingerprint string, certs []Cert) SlackSection {
	var sb strings.Builder

	// Subject
	sb.WriteString(">")
	sb.WriteString(certs[0].X509.Subject.CommonName)
	sb.WriteString("\n")

	// Expiry
	sb.WriteString("Expires in ")
	days := certs[0].ExpiresInDays()
	sb.WriteString(strconv.Itoa(days))
	if days == 1 {
		sb.WriteString(" day")
	} else {
		sb.WriteString(" days")
	}
	sb.WriteString("\n")

	// Secret keys
	if len(certs) == 1 {
		sb.WriteString("Secret:\n")
	} else {
		sb.WriteString("Secrets:\n")
	}
	for _, cert := range certs {
		sb.WriteString("`")
		sb.WriteString(cert.ToKey())
		sb.WriteString("`\n")
	}

	block := SlackSection{
		Type: "section",
		Text: SlackBlock{
			Type: "mrkdwn",
			Text: sb.String(),
		},
	}

	return block
}
