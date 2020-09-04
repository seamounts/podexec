package kube

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Config() (*restclient.Config, error) {
	// config, err := restclient.InClusterConfig()
	// if err == nil {
	// 	return config, nil
	// }

	kubeconfigPath := os.Getenv("KUBE_CONFIG")
	fmt.Println("-------kubeconfigPath", kubeconfigPath)
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// func ConfigForCluster(clusterID string) (*restclient.Config, error) {
// 	config, err := restclient.InClusterConfig()
// 	if err == nil {
// 		return config, nil
// 	}

// 	kubeconfigPath := os.Getenv("KUBE_CONFIG")
// 	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
// }

func Client() (*kubernetes.Clientset, error) {
	config, err := Config()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
