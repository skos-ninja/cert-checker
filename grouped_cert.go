package main

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"sort"
	"strconv"
	"strings"
)

type GroupedCert struct {
	Fingerprint   string
	ExpiresInDays int
	Certs         []Cert
}

func (g GroupedCert) ToSlackMessage() SlackSection {
	var sb strings.Builder

	// Subject
	sb.WriteString(">")
	sb.WriteString(g.Certs[0].X509.Subject.CommonName)
	sb.WriteString("\n")

	// Expiry
	days := g.ExpiresInDays
	if days < 0 {
		sb.WriteString("*EXPIRED ")
		days = int(math.Abs(float64(days)))
		sb.WriteString(strconv.Itoa(days))
		if days == 1 {
			sb.WriteString(" DAY")
		} else {
			sb.WriteString(" DAYS")
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
	if len(g.Certs) == 1 {
		sb.WriteString("Secret:\n")
	} else {
		sb.WriteString("Secrets:\n")
	}
	for _, cert := range g.Certs {
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

func groupCerts(certs []Cert) []GroupedCert {
	groupedCert := make(map[string][]Cert)
	for _, cert := range certs {
		cert := cert
		fingerprint := getMD5Hash(cert.X509.Raw)
		groupedCert[fingerprint] = append(groupedCert[fingerprint], cert)
	}

	gCerts := make([]GroupedCert, 0, len(groupedCert))
	for fingerprint, certs := range groupedCert {
		gCerts = append(gCerts, GroupedCert{
			Fingerprint:   fingerprint,
			ExpiresInDays: certs[0].ExpiresInDays(),
			Certs:         certs,
		})
	}

	// Sort our grouped into ascending expiry
	sort.Slice(gCerts, func(i, j int) bool {
		return gCerts[i].ExpiresInDays < gCerts[j].ExpiresInDays
	})

	return gCerts
}

func getMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
