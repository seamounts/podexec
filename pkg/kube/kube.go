package kube

import (
	"os"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Config() (*restclient.Config, error) {
	kubeconfigPath := os.Getenv("KUBE_CONFIG")
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func Client() (*kubernetes.Clientset, error) {
	config, err := Config()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
