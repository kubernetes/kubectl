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

	for i, arg := range os.Args[1:] {
		if arg != expectedArgs[i] {
			fmt.Printf("Expected arg %s, got arg %s", expectedArgs[i], arg)
			os.Exit(1)
		}
	}
	fmt.Println("Everything is fine")

	time.Sleep(10 * time.Minute)
}
