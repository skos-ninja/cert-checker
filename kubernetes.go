package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	allNamespaces = "*"
)

func getClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getNamespaces(clientSet *kubernetes.Clientset) (namespaces []string, err error) {
	nsInf := clientSet.CoreV1().Namespaces()
	opts := metav1.ListOptions{}
	ns, err := nsInf.List(opts)
	if err != nil {
		return nil, err
	}

	for _, namespace := range ns.Items {
		namespaces = append(namespaces, namespace.Name)
	}

	return namespaces, nil
}

func getSecrets(clientSet *kubernetes.Clientset, namespaces ...string) (secrets []corev1.Secret, err error) {
	secretItems := make([]corev1.Secret, 0)

	for _, namespace := range namespaces {
		secretInf := clientSet.CoreV1().Secrets(namespace)
		opts := metav1.ListOptions{}
		secrets, err := secretInf.List(opts)
		if err != nil {
			return nil, err
		}

		secretItems = append(secretItems, secrets.Items...)
	}

	return secretItems, nil
}
