package test

import (
	"fmt"
	"os/exec"
	"time"

	"net/url"

	"os"

	"io/ioutil"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"k8s.io/kubectl/pkg/framework/test/internal"
)

// APIServer knows how to run a kubernetes apiserver.
type APIServer struct {
	// URL is the address, a host and a port, the ApiServer should listen on for client connections.
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// Path is the path to the apiserver binary. If this is left as the empty
	// string, we will attempt to locate a binary, by checking for the
	// TEST_ASSET_KUBE_APISERVER environment variable, and the default test
	// assets directory.
	Path string

	// CertDir is a struct holding a path to a certificate directory and a function to cleanup that directory.
	CertDir       string
	actualCertDir string

	// Etcd is an implementation of a ControlPlaneProcess and is responsible to run Etcd and provide its coordinates.
	// If not specified, a brand new instance of Etcd is brought up.
	//
	// You can customise this if, e.g. you wish to use a already existing and running Etcd.
	// See the example `RemoteEtcd`.
	Etcd *Etcd

	// StopTimeout, StartTimeout specify the time the APIServer is allowed to take when stopping resp. starting
	// before and error is emitted.
	StopTimeout  time.Duration
	StartTimeout time.Duration

	session *gexec.Session
	stdOut  *gbytes.Buffer
	stdErr  *gbytes.Buffer
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

	etcdURLString := s.Etcd.URL.String()

	args := []string{
		"--authorization-mode=Node,RBAC",
		"--runtime-config=admissionregistration.k8s.io/v1alpha1",
		"--v=3", "--vmodule=",
		"--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota",
		"--admission-control-config-file=",
		"--bind-address=0.0.0.0",
		"--storage-backend=etcd3",
		fmt.Sprintf("--etcd-servers=%s", etcdURLString),
		fmt.Sprintf("--cert-dir=%s", s.actualCertDir),
		fmt.Sprintf("--insecure-port=%s", s.URL.Port()),
		fmt.Sprintf("--insecure-bind-address=%s", s.URL.Hostname()),
	}

	detectedStart := s.stdErr.Detect(fmt.Sprintf("Serving insecurely on %s", s.URL.Host))
	timedOut := time.After(s.StartTimeout)

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

func (s *APIServer) ensureInitialized() error {
	if s.Path == "" {
		s.Path = internal.BinPathFinder("kube-apiserver")
	}
	if s.URL == nil {
		am := &internal.AddressManager{}
		port, host, err := am.Initialize()
		if err != nil {
			return err
		}
		s.URL = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		}
	}
	if s.CertDir == "" {
		certDir, err := ioutil.TempDir("", "k8s_test_framework_")
		if err != nil {
			return err
		}
		s.actualCertDir = certDir
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

	if s.CertDir == "" {
		return os.RemoveAll(s.actualCertDir)
	}
	return nil
}
