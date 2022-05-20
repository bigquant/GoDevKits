package kube_utils

import (
	"context"
	"testing"

	"github.com/bigquant/GoDevKits/pkg/kube_utils"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateClient(t *testing.T) {
	opts := &kube_utils.KubeClientOpts{
		InCluster:      false,
		KubeConfigPath: "./kubeconfig",
		MasterURL:      "https://bigquant:6443",
		Namespace:      "bigquant",
		AllNamespaces:  false,
	}
	client := kube_utils.NewKubeClient(opts)
	kubeCli := client.GetClientSetOrDie()
	if res, err := kubeCli.CoreV1().Namespaces().List(
		context.Background(),
		metaV1.ListOptions{},
	); err != nil {
		t.Errorf("Cannot list namespaces")
	} else {
		t.Logf("Get res: %v", res.Items[0])
	}
}
