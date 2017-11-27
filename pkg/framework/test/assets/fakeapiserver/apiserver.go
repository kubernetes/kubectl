package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	expectedArgs := []string{
		"--authorization-mode=Node,RBAC",
		"--runtime-config=admissionregistration.k8s.io/v1alpha1",
		"--v=3", "--vmodule=",
		"--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota",
		"--admission-control-config-file=",
		"--bind-address=0.0.0.0",
		"--insecure-bind-address=127.0.0.1",
		"--insecure-port=8080",
		"--storage-backend=etcd3",
		"--etcd-servers=the etcd url",
	}
	numExpectedArgs := len(expectedArgs)
	numGivenArgs := len(os.Args) - 1

	if numGivenArgs < numExpectedArgs {
		fmt.Printf("Expected at least %d args, only got %d\n", numExpectedArgs, numGivenArgs)
		os.Exit(2)
	}

	for i, arg := range expectedArgs {
		givenArg := os.Args[i+1]
		if arg != givenArg {
			fmt.Printf("Expected arg %s, got arg %s\n", arg, givenArg)
			os.Exit(1)
		}
	}
	fmt.Println("Everything is fine")
	fmt.Fprintln(os.Stderr, "Serving insecurely on 127.0.0.1:8080")

	time.Sleep(10 * time.Minute)
}
