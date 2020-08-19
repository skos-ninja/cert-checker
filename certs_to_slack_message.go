package main

import (
	"math"
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
	days := certs[0].ExpiresInDays()
	if days < 0 {
		sb.WriteString("*EXPIRED ")
		days = int(math.Abs(float64(days)))
		sb.WriteString(strconv.Itoa(days))
		if days == 1 {
			sb.WriteString("DAY")
		} else {
			sb.WriteString("DAYS")
		}
		sb.WriteString(" AGO*")
	} else {
		sb.WriteString("Expires in ")
		sb.WriteString(strconv.Itoa(days))
		if days == 1 {
			sb.WriteString(" day")
		} else {
			sb.WriteString(" days")
		}
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
