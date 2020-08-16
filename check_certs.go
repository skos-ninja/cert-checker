package main

import (
	"time"
)

func checkCerts(certs []Cert) error {
	expiringCerts := make([]Cert, 0, len(certs))

	now := time.Now()
	expireAlertDate := now.Add(time.Duration(expiresInDays) * time.Duration(24) * time.Hour)
	for _, cert := range certs {
		cert := cert
		expiry := cert.X509.NotAfter

		if expiry.Before(expireAlertDate) {
			expiringCerts = append(expiringCerts, cert)
		}
	}

	if len(expiringCerts) > 0 {
		return alertExpiringCerts(expiringCerts)
	}

	return nil
}
