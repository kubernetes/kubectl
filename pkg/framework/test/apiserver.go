package test

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// APIServer knows how to run a kubernetes apiserver. Set it up with the path to a precompiled binary.
type APIServer struct {
	AddressManager AddressManager
	PathFinder     BinPathFinder
	ProcessStarter simpleSessionStarter
	CertDirManager certDirManager
	Etcd           ControlPlaneProcess
	StopTimeout    time.Duration
	StartTimeout   time.Duration
	session        SimpleSession
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
}

type certDirManager interface {
	Create() (string, error)
	Destroy() error
}

//go:generate counterfeiter . certDirManager

// URL returns the URL APIServer is listening on. Clients can use this to connect to APIServer.
func (s *APIServer) URL() (string, error) {
	if s.AddressManager == nil {
		return "", fmt.Errorf("APIServer's AddressManager is not initialized")
	}
	port, err := s.AddressManager.Port()
	if err != nil {
		return "", err
	}
	host, err := s.AddressManager.Host()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d", host, port), nil
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	s.ensureInitialized()

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

	etcdURLString, err := s.Etcd.URL()
	if err != nil {
		if etcdStopErr := s.Etcd.Stop(); etcdStopErr != nil {
			return fmt.Errorf("%s, %s", err.Error(), etcdStopErr.Error())
		}
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
		fmt.Sprintf("--etcd-servers=%s", etcdURLString),
		fmt.Sprintf("--cert-dir=%s", certDir),
		fmt.Sprintf("--insecure-port=%d", port),
		fmt.Sprintf("--insecure-bind-address=%s", addr),
	}

	detectedStart := s.stdErr.Detect(fmt.Sprintf("Serving insecurely on %s:%d", addr, port))
	timedOut := time.After(s.StartTimeout)

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

func (s *APIServer) ensureInitialized() {
	if s.PathFinder == nil {
		s.PathFinder = DefaultBinPathFinder
	}
	if s.AddressManager == nil {
		s.AddressManager = &DefaultAddressManager{}
	}
	if s.ProcessStarter == nil {
		s.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
			return gexec.Start(command, out, err)
		}
	}
	if s.CertDirManager == nil {
		s.CertDirManager = NewTempDirManager()
	}
	if s.Etcd == nil {
		s.Etcd = &Etcd{}
	}
	if s.StopTimeout == 0 {
		s.StopTimeout = 20 * time.Second
	}
	if s.StartTimeout == 0 {
		s.StartTimeout = 20 * time.Second
	}

	s.stdOut = gbytes.NewBuffer()
	s.stdErr = gbytes.NewBuffer()
}

// Stop stops this process gracefully, waits for its termination, and cleans up the cert directory.
func (s *APIServer) Stop() error {
	if s.session == nil {
		return nil
	}

	session := s.session.Terminate()
	detectedStop := session.Exited
	timedOut := time.After(s.StopTimeout)

	select {
	case <-detectedStop:
		break
	case <-timedOut:
		return fmt.Errorf("timeout waiting for apiserver to stop")
	}

	if err := s.Etcd.Stop(); err != nil {
		return err
	}

	return s.CertDirManager.Destroy()
}

// ExitCode returns the exit code of the process, if it has exited. If it hasn't exited yet, ExitCode returns -1.
func (s *APIServer) ExitCode() int {
	return s.session.ExitCode()
}

// Buffer implements the gbytes.BufferProvider interface and returns the stdout of the process
func (s *APIServer) Buffer() *gbytes.Buffer {
	return s.session.Buffer()
}
