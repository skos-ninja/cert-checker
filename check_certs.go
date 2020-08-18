package main

import (
	"github.com/sirupsen/logrus"
)

func checkCerts(certs []Cert) error {
	expiringCerts := make([]Cert, 0, len(certs))

	for _, cert := range certs {
		cert := cert
		certExpires := cert.ExpiresInDays()

		logrus.WithFields(logrus.Fields{
			"namespace": cert.Namespace,
			"name":      cert.Name,
			"key":       cert.Key,
			"subject":   cert.X509.Subject,
		}).Infof("Expires in %v", certExpires)

		if certExpires <= expiresInDays {
			expiringCerts = append(expiringCerts, cert)
		}
	}

	if len(expiringCerts) > 0 {
		return alertExpiringCerts(expiringCerts)
	}

	return nil
}
