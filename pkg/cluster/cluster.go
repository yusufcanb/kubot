package cluster

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type Cluster struct {
	defaultNamespace string

	client *kubernetes.Clientset
	config *rest.Config
}

func (c *Cluster) DefaultNamespace() string {
	return c.defaultNamespace
}

func (c *Cluster) Client() *kubernetes.Clientset {
	return c.client
}

func (c *Cluster) Config() *rest.Config {
	return c.config
}

func NewCluster(kubeConfigPath string, namespace string) (*Cluster, error) {
	cluster := Cluster{}
	cluster.defaultNamespace = namespace

	if home, _ := os.UserHomeDir(); home != "" {
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	cluster.config = config

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cluster.client = client

	return &cluster, nil
}
