package test

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// APIServer knows how to run a kubernetes apiserver. Set it up with the path to a precompiled binary.
type APIServer struct {
	AddressManager AddressManager
	PathFinder     BinPathFinder
	ProcessStarter simpleSessionStarter
	CertDirManager certDirManager
	Config         *APIServerConfig
	Etcd           FixtureProcess
	session        SimpleSession
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
}

type certDirManager interface {
	Create() (string, error)
	Destroy() error
}

//go:generate counterfeiter . certDirManager

// NewAPIServer creates a new APIServer Fixture Process
func NewAPIServer(config *APIServerConfig) (*APIServer, error) {
	starter := func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
		return gexec.Start(command, out, err)
	}

	etcd, err := NewEtcd()
	if err != nil {
		return nil, err
	}

	return &APIServer{
		Config:         config,
		ProcessStarter: starter,
		CertDirManager: NewTempDirManager(),
		Etcd:           etcd,
	}, nil
}

// URL returns the URL APIServer is listening on. Clients can use this to connect to APIServer.
func (s *APIServer) URL() string {
	// TODO handle errors
	port, _ := s.AddressManager.Port()
	host, _ := s.AddressManager.Host()
	return fmt.Sprintf("http://%s:%d", host, port)
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	if s.PathFinder == nil {
		s.PathFinder = DefaultBinPathFinder
	}
	if s.AddressManager == nil {
		s.AddressManager = &DefaultAddressManager{}
	}
	if err := s.Config.Validate(); err != nil {
		return err
	}

	port, addr, err := s.AddressManager.Initialize("localhost")
	if err != nil {
		return err
	}

	certDir, err := s.CertDirManager.Create()
	if err != nil {
		return err
	}

	err = s.Etcd.Start()
	if err != nil {
		return err
	}

	s.stdOut = gbytes.NewBuffer()
	s.stdErr = gbytes.NewBuffer()

	args := []string{
		"--authorization-mode=Node,RBAC",
		"--runtime-config=admissionregistration.k8s.io/v1alpha1",
		"--v=3", "--vmodule=",
		"--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota",
		"--admission-control-config-file=",
		"--bind-address=0.0.0.0",
		"--storage-backend=etcd3",
		fmt.Sprintf("--etcd-servers=%s", s.Etcd.URL()),
		fmt.Sprintf("--cert-dir=%s", certDir),
		fmt.Sprintf("--insecure-port=%d", port),
		fmt.Sprintf("--insecure-bind-address=%s", addr),
	}

	detectedStart := s.stdErr.Detect(fmt.Sprintf("Serving insecurely on %s:%d", addr, port))
	timedOut := time.After(20 * time.Second)

	command := exec.Command(s.PathFinder("kube-apiserver"), args...)
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
		s.Etcd.Stop()
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
