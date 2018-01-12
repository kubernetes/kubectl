package test

import (
	"fmt"
	"time"

	"net/url"

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
	CertDir string

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

	processState *internal.ProcessState
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	var err error

	if s.Etcd == nil {
		s.Etcd = &Etcd{}
	}

	err = s.Etcd.Start()
	if err != nil {
		return err
	}

	s.processState, err = internal.NewProcessState(
		"kube-apiserver",
		s.Path,
		s.URL,
		s.CertDir,
		s.StartTimeout, s.StopTimeout,
	)
	if err != nil {
		return err
	}

	s.processState.Args = []string{
		"--authorization-mode=Node,RBAC",
		"--runtime-config=admissionregistration.k8s.io/v1alpha1",
		"--v=3", "--vmodule=",
		"--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota",
		"--admission-control-config-file=",
		"--bind-address=0.0.0.0",
		"--storage-backend=etcd3",
		fmt.Sprintf("--etcd-servers=%s", s.Etcd.processState.URL.String()),
		fmt.Sprintf("--cert-dir=%s", s.processState.Dir),
		fmt.Sprintf("--insecure-port=%s", s.processState.URL.Port()),
		fmt.Sprintf("--insecure-bind-address=%s", s.processState.URL.Hostname()),
	}
	if err != nil {
		return err
	}

	s.processState.Start(fmt.Sprintf("Serving insecurely on %s", s.processState.URL.Host))

	return err
}

// Stop stops this process gracefully, waits for its termination, and cleans up the cert directory.
func (s *APIServer) Stop() error {
	err := s.processState.Stop()
	if err != nil {
		return err
	}

	return s.Etcd.Stop()
}
