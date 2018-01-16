package test

import (
	"fmt"
	"net/url"
	"time"

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
	EtcdURL *url.URL

	// StopTimeout, StartTimeout specify the time the APIServer is allowed to take when stopping resp. starting
	// before and error is emitted.
	StopTimeout  time.Duration
	StartTimeout time.Duration

	processState *internal.ProcessState
}

// Start starts the apiserver, waits for it to come up, and returns an error, if occoured.
func (s *APIServer) Start() error {
	var err error

	s.processState = &internal.ProcessState{}

	s.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-apiserver",
		s.URL,
		s.CertDir,
		s.Path,
		s.StartTimeout,
		s.StopTimeout,
	)
	if err != nil {
		return err
	}

	s.processState.Args, err = internal.MakeAPIServerArgs(
		s.processState.DefaultedProcessInput,
		s.EtcdURL,
	)
	if err != nil {
		return err
	}

	s.processState.StartMessage = fmt.Sprintf(
		"Serving insecurely on %s",
		s.processState.URL.Host,
	)

	return s.processState.Start()
}

// Stop stops this process gracefully, waits for its termination, and cleans up the cert directory.
func (s *APIServer) Stop() error {
	return s.processState.Stop()
}
