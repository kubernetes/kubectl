package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	expectedArgs := []*regexp.Regexp{
		regexp.MustCompile("^--authorization-mode=Node,RBAC$"),
		regexp.MustCompile("^--runtime-config=admissionregistration.k8s.io/v1alpha1$"),
		regexp.MustCompile("^--v=3$"),
		regexp.MustCompile("^--vmodule=$"),
		regexp.MustCompile("^--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota$"),
		regexp.MustCompile("^--admission-control-config-file=$"),
		regexp.MustCompile("^--bind-address=0.0.0.0$"),
		regexp.MustCompile("^--insecure-bind-address=127.0.0.1$"),
		regexp.MustCompile("^--insecure-port=8080$"),
		regexp.MustCompile("^--storage-backend=etcd3$"),
		regexp.MustCompile("^--etcd-servers=the etcd url$"),
		regexp.MustCompile("^--cert-dir=.*"),
	}
	numExpectedArgs := len(expectedArgs)
	numGivenArgs := len(os.Args) - 1

	if numGivenArgs < numExpectedArgs {
		fmt.Printf("Expected at least %d args, only got %d\n", numExpectedArgs, numGivenArgs)
		os.Exit(2)
	}

	for i, argRegexp := range expectedArgs {
		givenArg := os.Args[i+1]
		if !argRegexp.MatchString(givenArg) {
			fmt.Printf("Expected arg '%s' to match '%s'\n", givenArg, argRegexp.String())
			os.Exit(1)
		}
	}
	fmt.Println("Everything is fine")
	fmt.Fprintln(os.Stderr, "Serving insecurely on 127.0.0.1:8080")

	time.Sleep(10 * time.Minute)
}
