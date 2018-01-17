package test

import (
	"fmt"
	"net/url"
	"time"

	"k8s.io/kubectl/pkg/framework/test/internal"
)

// APIServer knows how to run a kubernetes apiserver.
type APIServer struct {
	// URL is the address the ApiServer should listen on for client connections.
	//
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// Path is the path to the apiserver binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_KUBE_APISERVER environment variable, and
	// the default test assets directory.
	Path string

	// CertDir is a path to a directory containing whatever certificates the
	// APIServer will need.
	//
	// If left unspecified, then the Start() method will create a fresh temporary
	// directory, and the Stop() method will clean it up.
	CertDir string

	// EtcdURL is the URL of the Etcd the APIServer should use.
	//
	// If this is not specified, the Start() method will return an error.
	EtcdURL *url.URL

	// StartTimeout, StopTimeout specify the time the APIServer is allowed to
	// take when starting and stoppping before an error is emitted.
	//
	// If not specified, these default to 20 seconds.
	StartTimeout time.Duration
	StopTimeout  time.Duration

	processState *internal.ProcessState
}

// Start starts the apiserver, waits for it to come up, and returns an error,
// if occurred.
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

// Stop stops this process gracefully, waits for its termination, and cleans up
// the CertDir if necessary.
func (s *APIServer) Stop() error {
	return s.processState.Stop()
}
