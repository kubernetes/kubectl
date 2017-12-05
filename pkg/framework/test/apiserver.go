package test

import (
	"fmt"
	"os/exec"
	"time"

	"io"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// APIServer knows how to run a kubernetes apiserver. Set it up with the path to a precompiled binary.
type APIServer struct {
	// The path to the apiserver binary
	Path           string
	EtcdURL        string
	ProcessStarter simpleSessionStarter
	CertDirManager certDirManager
	session        SimpleSession
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
}

type certDirManager interface {
	Create() (string, error)
	Destroy() error
}

//go:generate counterfeiter . certDirManager

func NewAPIServer(pathToAPIServer, etcdURL string) *APIServer {
	starter := func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
		return gexec.Start(command, out, err)
	}

	apiserver := &APIServer{
		Path:           pathToAPIServer,
		EtcdURL:        etcdURL,
		ProcessStarter: starter,
		CertDirManager: NewTempDirManager(),
	}

	return apiserver
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	s.stdOut = gbytes.NewBuffer()
	s.stdErr = gbytes.NewBuffer()

	certDir, err := s.CertDirManager.Create()
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
	s.session, err = s.ProcessStarter(command, s.stdOut, s.stdErr)
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

// Stop stops this process gracefully, waits for its termination, and cleans up the cert directory.
func (s *APIServer) Stop() {
	if s.session != nil {
		s.session.Terminate()
		s.session.Wait(20 * time.Second)
		err := s.CertDirManager.Destroy()
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
