package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

type test struct {
	name string
	args args
	want bool
}

type args struct {
	cert                      Cert
	minCert, expires, expired int
	time                      time.Time
}

func generateCert(t *testing.T, begin time.Time, expiry time.Time) Cert {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Testing Ca"},
			CommonName:   "Testing Cert",
		},
		NotBefore: begin,
		NotAfter:  expiry,

		KeyUsage:              0,
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		BasicConstraintsValid: true,
	}

	return Cert{
		Name:      "test",
		Namespace: "test",
		Key:       "test",

		X509: cert,
	}
}

func shortLengthCert(t *testing.T) test {
	begin := time.Now().Add(-1 * time.Hour)
	expiry := begin.Add(1 * time.Hour)

	return test{
		name: "Short length cert",
		args: args{
			cert:    generateCert(t, begin, expiry),
			minCert: 1,
			expires: 30,
			expired: 90,
			time:    begin,
		},
		want: false,
	}
}

func expiredCert(t *testing.T) test {
	now := time.Now()
	begin := now.Add(-96 * time.Hour) // Cert was issued 96 hours ago
	expiry := now.Add(-1 * time.Hour) // Cert expired 1 hour ago

	return test{
		name: "Expired cert",
		args: args{
			cert:    generateCert(t, begin, expiry),
			minCert: 0,
			expires: 30,
			expired: 90,
			time:    now,
		},
		want: true,
	}
}

func expiringSoonCert(t *testing.T) test {
	now := time.Now()
	begin := now.Add(-96 * time.Hour) // Cert was issued 96 hours ago
	expiry := now.Add(24 * time.Hour) // Cert expires in 24 hours

	return test{
		name: "Expires in cert",
		args: args{
			cert:    generateCert(t, begin, expiry),
			minCert: 0,
			expires: 30,
			expired: 90,
			time:    now,
		},
		want: true,
	}
}

func oldExpiredCert(t *testing.T) test {
	now := time.Now()
	begin := now.Add(-(180 * (24 * time.Hour)))  // Cert was issued 180 days ago
	expiry := now.Add(-(120 * (24 * time.Hour))) // Cert expired 120 days ago

	return test{
		name: "Old expired cert",
		args: args{
			cert:    generateCert(t, begin, expiry),
			minCert: 0,
			expires: 30,
			expired: 90,
			time:    now,
		},
		want: false,
	}
}

func validCert(t *testing.T) test {
	now := time.Now()
	begin := now.Add(-(30 * (24 * time.Hour))) // Cert was issued 30 days ago
	expiry := now.Add(60 * (24 * time.Hour))   // Cert expires in 30 days

	return test{
		name: "Valid cert",
		args: args{
			cert:    generateCert(t, begin, expiry),
			minCert: 0,
			expires: 30,
			expired: 90,
			time:    now,
		},
		want: false,
	}
}

func Test_isCertFlagged(t *testing.T) {
	tests := []test{
		shortLengthCert(t),
		expiredCert(t),
		expiringSoonCert(t),
		oldExpiredCert(t),
		validCert(t),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got :=
				isCertFlagged(
					tt.args.minCert,
					tt.args.expires,
					tt.args.expired,
					tt.args.cert,
					tt.args.time); got != tt.want {
				t.Errorf("isCertFlagged() = %v, want %v", got, tt.want)
			}
		})
	}
}
