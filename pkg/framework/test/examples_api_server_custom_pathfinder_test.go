package test_test

import (
	"fmt"

	. "k8s.io/kubectl/pkg/framework/test"
)

func ExampleAPIServer_specialPathFinder() {
	apiServer := &APIServer{
		PathFinder: myCustomPathFinder,
	}

	fmt.Println(apiServer)
}

func myCustomPathFinder(_ string) string {
	return "/usr/local/bin/kube/1.21-alpha1/special-kube-api"
}
