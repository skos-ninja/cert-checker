package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

var (
	errSecretNotPEM  = errors.New("secret is not a valid PEM")
	errSecretNotCert = errors.New("secret is not a valid certificate")
)

func parseSecret(secret v1.Secret) ([]Cert, error) {
	certs := make([]Cert, 0)
	logger := logrus.WithFields(logrus.Fields{
		"namspeace": secret.Namespace,
		"name":      secret.Name,
	})

	for key, data := range secret.Data {
		logger := logger.WithField("key", key)
		cert, err := parseCert(data)
		switch err {
		case errSecretNotPEM:
			logger.Debug("Skipping secret as not PEM")
			continue
		case errSecretNotCert:
			logger.Info("Skipping secret as not a valid certificate")
			continue
		default:
			logger.WithError(err).Error("Failed to parse secret")
		}

		certs = append(certs, Cert{
			Namespace: secret.Namespace,
			Name:      secret.Name,
			Key:       key,

			X509: cert,
		})
	}

	return certs, nil
}

func parseCert(data []byte) (*x509.Certificate, error) {
	// Try to base64 parse our cert before pem decoding
	base64String := string(data)
	parsedData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		parsedData = data
	}

	block, _ := pem.Decode(parsedData)
	if block == nil {
		// Data isn't a valid pem
		return nil, errSecretNotPEM
	}

	// Ensure we ignore secrets that aren't a certificate
	if block.Type != "CERTIFICATE" {
		return nil, errSecretNotCert
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		// PEM isn't a valid certificate
		log.Printf("Error parsing certificate %s\n", err.Error())
		return nil, err
	}

	return cert, nil
}
