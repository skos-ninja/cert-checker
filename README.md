# cert-checker

cert-checker allows you to scan for certificates in a kubernetes cluster and alert you via slack when they are set to expire soon.

![example message](https://github.com/skos-ninja/cert-checker/blob/master/example/example-msg.png?raw=true)

## Parameters

- `--expires-in-days` sets the number of days before the certificates expire to alert
- `--max-expired-in-days` set the number of days after the certificates expire to stop alerting
- `--min-cert-length-in-days` set the minimum number of days a certificate has to be valid for to alert
- `--slack-webhook` defines the webhook URL of the slack app
- `--namespace` can be passed multiple times to define each namespace that the app should scan (for all namespaces provide `*`)
- `--environment` specifies the environment context for the slack message
- `--alert-no-flagged-certs` enables a slack message for when no certificates are flagged

## Deployment

A slack app must be created for you to deploy this app. You can follow the steps [here](https://api.slack.com/messaging/webhooks#getting_started).

An example of a kubernetes deployment has been provided in [deployment/kubernetes.yml](deployment/kubernetes.yml) which provides a namespace, service account and cronjob.