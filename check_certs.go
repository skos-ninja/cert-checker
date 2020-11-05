package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

func checkCerts(certs []Cert, t time.Time) []Cert {
	expiringCerts := make([]Cert, 0, len(certs))

	for _, cert := range certs {
		cert := cert

		flagged := isCertFlagged(minCertLengthInDays, expiresInDays, maxExpiredInDays, cert, t)
		if flagged {
			expiringCerts = append(expiringCerts, cert)
		}
	}

	return expiringCerts
}

func isCertFlagged(minCert, expires, expired int, cert Cert, t time.Time) bool {
	log := logrus.WithFields(logrus.Fields{
		"namespace": cert.Namespace,
		"name":      cert.Name,
		"key":       cert.Key,
		"subject":   cert.X509.Subject,
	})

	certExpires := cert.ExpiresInDays(t)
	log.Infof("Expires in %v", certExpires)

	validPeriod := cert.ValidPeriodInDays()
	if validPeriod < minCert {
		log.Infof("Cert valid period too short: %v", validPeriod)
		return false
	}

	if certExpires <= expires {
		log.Infof("Cert expires within: %v (%v)", certExpires, expires)

		if certExpires >= -(expired) {
			return true
		}
		log.Infof("Cert expired for too long: %v (%v)", certExpires, expired)
	}

	return false
}
