package test

import (
	"os/exec"

	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// APIServer knows how to run a kubernetes apiserver. Set it up with the path to a precompiled binary.
type APIServer struct {
	// The path to the apiserver binary
	Path    string
	session *gexec.Session
}

// Start starts the apiserver, and returns a gexec.Session. To stop it again, call Terminate and Wait on that session.
func (s *APIServer) Start(etcdURL string) error {
	args := []string{
		"--authorization-mode=Node,RBAC",
		"--runtime-config=admissionregistration.k8s.io/v1alpha1",
		"--v=3", "--vmodule=",
		"--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota",
		"--admission-control-config-file=",
		"--bind-address=0.0.0.0",
		"--insecure-bind-address=127.0.0.1",
		"--insecure-port=8080",
		"--storage-backend=etcd3",
		fmt.Sprintf("--etcd-servers=%s", etcdURL),
	}

	command := exec.Command(s.Path, args...)
	var err error
	s.session, err = gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	return err
}

// Stop stops this process gracefully.
func (s *APIServer) Stop() {
	if s.session != nil {
		s.session.Terminate()
	}
}

// ExitCode returns the exit code of the process, if it has exited. If it hasn't exited yet, ExitCode returns -1.
func (s *APIServer) ExitCode() int {
	return s.session.ExitCode()
}

// Buffer implements the gbytes.BufferProvider interface and returns the stdout of the process
func (s *APIServer) Buffer() *gbytes.Buffer {
	return s.session.Buffer()
}
