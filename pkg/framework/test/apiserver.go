package test

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"net/url"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"k8s.io/kubectl/pkg/framework/test/internal"
)

// APIServer knows how to run a kubernetes apiserver.
type APIServer struct {
	// Address is the address, a host and a port, the ApiServer should listen on for client connections.
	// If this is not specified, we default to a random free port on localhost.
	Address *url.URL

	// Path is the path to the apiserver binary. If this is left as the empty
	// string, we will attempt to locate a binary, by checking for the
	// TEST_ASSET_KUBE_APISERVER environment variable, and the default test
	// assets directory.
	Path string

	// ProcessStarter is a way to hook into how a the APIServer process is started. By default `gexec.Start(...)` is
	// used to run the process.
	//
	// You can customize this if, e.g. you want to pass additional arguments or do extra logging.
	// See the `SpecialPathFinder` example.
	ProcessStarter SimpleSessionStarter

	// CertDir is a struct holding a path to a certificate directory and a function to cleanup that directory.
	CertDir *CleanableDirectory

	// Etcd is an implementation of a ControlPlaneProcess and is responsible to run Etcd and provide its coordinates.
	// If not specified, a brand new instance of Etcd is brought up.
	//
	// You can customise this if, e.g. you wish to use a already existing and running Etcd.
	// See the example `RemoteEtcd`.
	Etcd ControlPlaneProcess

	// StopTimeout, StartTimeout specify the time the APIServer is allowed to take when stopping resp. starting
	// before and error is emitted.
	StopTimeout  time.Duration
	StartTimeout time.Duration

	session SimpleSession
	stdOut  *gbytes.Buffer
	stdErr  *gbytes.Buffer
}

// URL returns the URL APIServer is listening on. Clients can use this to connect to APIServer.
func (s *APIServer) URL() (string, error) {
	if s.Address == nil {
		return "", fmt.Errorf("APIServer's Address is not initialized or configured")
	}
	return s.Address.String(), nil
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	err := s.ensureInitialized()
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
		fmt.Sprintf("--cert-dir=%s", s.CertDir.Path),
		fmt.Sprintf("--insecure-port=%s", s.Address.Port()),
		fmt.Sprintf("--insecure-bind-address=%s", s.Address.Hostname()),
	}

	detectedStart := s.stdErr.Detect(fmt.Sprintf("Serving insecurely on %s", s.Address.Host))
	timedOut := time.After(s.StartTimeout)

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

func (s *APIServer) ensureInitialized() error {
	if s.Path == "" {
		s.Path = internal.BinPathFinder("kube-apiserver")
	}
	if s.Address == nil {
		am := &internal.AddressManager{}
		port, host, err := am.Initialize()
		if err != nil {
			return err
		}
		s.Address = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		}
	}
	if s.ProcessStarter == nil {
		s.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
			return gexec.Start(command, out, err)
		}
	}
	if s.CertDir == nil {
		certDir, err := newDirectory()
		if err != nil {
			return err
		}
		s.CertDir = certDir
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

	return nil
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

	if s.CertDir.Cleanup == nil {
		return nil
	}
	return s.CertDir.Cleanup()
}

// ExitCode returns the exit code of the process, if it has exited. If it hasn't exited yet, ExitCode returns -1.
func (s *APIServer) ExitCode() int {
	return s.session.ExitCode()
}

// Buffer implements the gbytes.BufferProvider interface and returns the stdout of the process
func (s *APIServer) Buffer() *gbytes.Buffer {
	return s.session.Buffer()
}
