package test

import (
	"fmt"
	"io"
	"net/url"
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
	ProcessStarter simpleSessionStarter
	CertDirManager certDirManager
	Config         *APIServerConfig
	session        SimpleSession
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
}

type certDirManager interface {
	Create() (string, error)
	Destroy() error
}

//go:generate counterfeiter . certDirManager

var apiServerBinPathFinder = DefaultBinPathFinder

// NewAPIServer creates a new APIServer Fixture Process
func NewAPIServer(config *APIServerConfig) *APIServer {
	starter := func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
		return gexec.Start(command, out, err)
	}

	return &APIServer{
		Path:           apiServerBinPathFinder("kube-apiserver"),
		Config:         config,
		ProcessStarter: starter,
		CertDirManager: NewTempDirManager(),
	}
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	if err := s.Config.Validate(); err != nil {
		return err
	}

	s.stdOut = gbytes.NewBuffer()
	s.stdErr = gbytes.NewBuffer()

	certDir, err := s.CertDirManager.Create()
	if err != nil {
		return err
	}

	clientURL, err := url.Parse(s.Config.APIServerURL)
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
		"--storage-backend=etcd3",
		fmt.Sprintf("--etcd-servers=%s", s.Config.EtcdURL),
		fmt.Sprintf("--cert-dir=%s", certDir),
		fmt.Sprintf("--insecure-port=%s", clientURL.Port()),
		fmt.Sprintf("--insecure-bind-address=%s", clientURL.Hostname()),
	}

	detectedStart := s.stdErr.Detect(fmt.Sprintf("Serving insecurely on %s", clientURL.Host))
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
