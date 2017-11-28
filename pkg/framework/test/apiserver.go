package test

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// APIServer knows how to run a kubernetes apiserver. Set it up with the path to a precompiled binary.
type APIServer struct {
	// The path to the apiserver binary
	Path           string
	EtcdURL        string
	session        *gexec.Session
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
	certDirManager certDirManager
}

type certDirManager interface {
	Create() (string, error)
	Destroy() error
}

// Start starts the apiserver, and returns a gexec.Session. To stop it again, call Terminate and Wait on that session.
func (s *APIServer) Start() error {
	s.certDirManager = NewTempDirManager()
	s.stdOut = gbytes.NewBuffer()
	s.stdErr = gbytes.NewBuffer()

	certDir, err := s.certDirManager.Create()
	if err != nil {
		return err
	}

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
		fmt.Sprintf("--etcd-servers=%s", s.EtcdURL),
		fmt.Sprintf("--cert-dir=%s", certDir),
	}

	detectedStart := s.stdErr.Detect("Serving insecurely on 127.0.0.1:8080")
	timedOut := time.After(20 * time.Second)

	command := exec.Command(s.Path, args...)
	s.session, err = gexec.Start(command, s.stdOut, s.stdErr)
	if err != nil {
		return err
	}

	select {
	case <-detectedStart:
		return nil
	case <-timedOut:
		return fmt.Errorf("timeout waiting for apiserver to start serving")
	}
}

// Stop stops this process gracefully.
func (s *APIServer) Stop() {
	if s.session != nil {
		s.session.Terminate().Wait(20 * time.Second)
		err := s.certDirManager.Destroy()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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
