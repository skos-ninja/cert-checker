package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:  "cert-checker",
		Args: cobra.ExactArgs(0),
		RunE: runE,
	}

	namespaces          = []string{"default"}
	expiresInDays       = 31
	maxExpiredInDays    = 31
	minCertLengthInDays = 0
	slackWebHook        = ""
	environmentString   = ""
)

func main() {
	cmd.Flags().StringArrayVarP(&namespaces, "namespace", "n", namespaces, "Define a namespace to scan")
	cmd.Flags().IntVarP(&expiresInDays, "expires-in-days", "e", expiresInDays, "Sets the number of days before expiry to alert")
	cmd.Flags().IntVarP(&maxExpiredInDays, "max-expired-in-days", "m", maxExpiredInDays, "Sets the number of days after expiry to stop alerting")
	cmd.Flags().IntVarP(&minCertLengthInDays, "min-cert-length-in-days", "l", minCertLengthInDays, "Sets the minimum number of days the certificate has to be valid before it is considered for alerting")
	cmd.Flags().StringVarP(&slackWebHook, "slack-webhook", "s", slackWebHook, "Slack webhook url for the client to alert with")
	cmd.Flags().StringVar(&environmentString, "environment", environmentString, "Adds an environment to your expired certs message")

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func runE(cmd *cobra.Command, args []string) error {
	now := time.Now()

	client, err := getClientSet()
	if err != nil {
		return err
	}

	if len(namespaces) == 1 && namespaces[0] == allNamespaces {
		namespaces, err = getNamespaces(client)
		if err != nil {
			return err
		}
	}

	secrets, err := getSecrets(client, namespaces...)
	if err != nil {
		return err
	}

	certs := make([]Cert, 0)
	for _, secret := range secrets {
		c, err := parseSecret(secret)
		if err != nil {
			return err
		}

		certs = append(certs, c...)
	}

	return checkCerts(certs, now)
}
