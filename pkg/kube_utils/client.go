package kube_utils

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdApi "k8s.io/client-go/tools/clientcmd/api"
)

// KubeClientOpts: kubernetes client options
type KubeClientOpts struct {
	InCluster      bool
	KubeConfigPath string // when not in cluster
	MasterURL      string // when not in cluster
	Namespace      string
	AllNamespaces  bool
}

type KubeClient interface {
	GetConfig() (*rest.Config, error)
	GetConfigOrDie() *rest.Config
	Namespace() string
	DefaultNamespace() string
	GetClientSet() (*kubernetes.Clientset, error)
	GetClientSetOrDie() *kubernetes.Clientset
}

// NewKubeClient func
func NewKubeClient(opts *KubeClientOpts) KubeClient {
	return &kubeClient{KubeClientOpts: *opts}
}

type kubeClient struct {
	KubeClientOpts
	clientConfigOnce sync.Once
	clientConfig     *rest.Config
	clientSetOnce    sync.Once
	clientSet        *kubernetes.Clientset
}

func (c *kubeClient) GetConfig() (*rest.Config, error) {
	var err error
	c.clientConfigOnce.Do(func() {
		if c.InCluster {
			c.clientConfig, err = rest.InClusterConfig()
		} else {
			loadingConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{
					ExplicitPath: c.KubeConfigPath,
				},
				&clientcmd.ConfigOverrides{
					ClusterInfo: clientcmdApi.Cluster{
						Server: c.MasterURL,
					},
				},
			)
			c.clientConfig, err = loadingConfig.ClientConfig()
		}
	})
	return c.clientConfig, err
}

func (c *kubeClient) GetConfigOrDie() *rest.Config {
	config, err := c.GetConfig()
	if err != nil {
		log.Printf("Failed to get kubernetes config: %v\n", err)
		panic(err)
	}
	return config
}

func (c *kubeClient) Namespace() string {
	if c.AllNamespaces {
		return metaV1.NamespaceAll
	}
	if !c.AllNamespaces && c.KubeClientOpts.Namespace == "" {
		c.KubeClientOpts.Namespace = c.DefaultNamespace()
		return c.KubeClientOpts.Namespace
	}
	return c.KubeClientOpts.Namespace
}

const (
	SA_NS_FILE_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func (c *kubeClient) DefaultNamespace() string {
	if data, err := ioutil.ReadFile(SA_NS_FILE_PATH); err == nil {
		if ns := strings.TrimSpace(string(data)); ns != "" {
			return ns
		}
	}
	return "default"
}

func (c *kubeClient) GetClientSet() (*kubernetes.Clientset, error) {
	var err error
	c.clientSetOnce.Do(func() {
		c.clientSet, err = kubernetes.NewForConfig(c.GetConfigOrDie())
	})
	return c.clientSet, err
}

func (c *kubeClient) GetClientSetOrDie() *kubernetes.Clientset {
	clientSet, err := c.GetClientSet()
	if err != nil {
		log.Printf("Failed to get kubernetes client set: %v\n", err)
		panic(err)
	}
	return clientSet
}
