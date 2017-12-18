package test_test

import (
	"io"
	"os/exec"

	"fmt"

	"github.com/onsi/gomega/gexec"
	. "k8s.io/kubectl/pkg/framework/test"
)

func ExampleAPIServer_forcedArgsProcessStarter() {
	apiServer := &APIServer{
		ProcessStarter: forceArgsStarter,
	}
	fmt.Println(apiServer)
}

func forceArgsStarter(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
	forcedArgs := []string{
		"--target-ram-mb=1024",
		"--allow-privileged=false",
		"--authorization-mode=RBAC,AlwaysDeny",
	}

	newArgs := append(command.Args, forcedArgs...)

	command.Args = newArgs

	return gexec.Start(command, out, err)
}
