package tools

import (
	// for kubernetes crd
	_ "k8s.io/code-generator"
	_ "k8s.io/code-generator/cmd/client-gen"
	_ "k8s.io/code-generator/cmd/deepcopy-gen"
	_ "k8s.io/code-generator/cmd/defaulter-gen"
	_ "k8s.io/code-generator/cmd/informer-gen"

	// for kubernetes apiextensions
	_ "k8s.io/kube-openapi/cmd/openapi-gen"

	// kubebuilder
	_ "sigs.k8s.io/kubebuilder/v3/cmd"
)
