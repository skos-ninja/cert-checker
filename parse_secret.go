package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"

	v1 "k8s.io/api/core/v1"
)

var (
	ErrSecretNotCert = errors.New("secret is not a valid certificate")
)

func parseSecret(secret v1.Secret) ([]Cert, error) {
	certs := make([]Cert, 0)

	for key, data := range secret.Data {
		cert, err := parseCert(data)
		if err == ErrSecretNotCert {
			continue
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
	certPEM := string(data)
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		// Data isn't a valid pem
		return nil, ErrSecretNotCert
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		// PEM isn't a valid certificate
		// This is probably because it's a key however this should probably be checked 🤔
		log.Printf("Error parsing certificate %s\n", err.Error())
		return nil, ErrSecretNotCert
	}

	return cert, nil
}
