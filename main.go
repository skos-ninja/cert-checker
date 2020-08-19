package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:  "cert-checker",
		Args: cobra.ExactArgs(0),
		RunE: runE,
	}

	namespaces    = []string{"default"}
	expiresInDays = 31
	slackWebHook  = ""
)

func main() {
	cmd.Flags().StringArrayVarP(&namespaces, "namespace", "n", namespaces, "Define a namespace to scan")
	cmd.Flags().IntVarP(&expiresInDays, "expires-in-days", "e", expiresInDays, "Sets the number of days before expiry to alert")
	cmd.Flags().StringVarP(&slackWebHook, "slack-webhook", "s", slackWebHook, "Slack webhook url for the client to alert with")

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func runE(cmd *cobra.Command, args []string) error {
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

	return checkCerts(certs)
}
