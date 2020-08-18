package main

import (
	"crypto/x509"
	"fmt"
	"math"
	"time"
)

type Cert struct {
	Namespace string
	Name      string
	Key       string

	X509 *x509.Certificate
}

func (c Cert) ToKey() string {
	return fmt.Sprintf("%s/%s/%s", c.Namespace, c.Name, c.Key)
}

func (c Cert) ExpiresInDays() int {
	expiry := c.X509.NotAfter
	expiresInDays := expiry.Sub(time.Now()).Hours() / 24

	return int(math.Floor(expiresInDays))
}
