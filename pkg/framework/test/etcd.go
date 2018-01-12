package test

import (
	"fmt"
	"time"

	"net/url"

	"k8s.io/kubectl/pkg/framework/test/internal"
)

// Etcd knows how to run an etcd server.
//
// The documentation and examples for the Etcd's properties can be found in
// in the documentation for the `APIServer`, as both implement a `ControlPaneProcess`.
type Etcd struct {
	URL          *url.URL
	Path         string
	DataDir      string
	StopTimeout  time.Duration
	StartTimeout time.Duration

	processState *internal.ProcessState
}

// Start starts the etcd, waits for it to come up, and returns an error, if occoured.
func (e *Etcd) Start() error {
	var err error

	e.processState, err = internal.NewProcessState(
		"etcd",
		e.Path,
		e.URL,
		e.DataDir,
		e.StartTimeout, e.StopTimeout,
	)
	if err != nil {
		return err
	}

	e.processState.Args = []string{
		"--debug",
		"--listen-peer-urls=http://localhost:0",
		fmt.Sprintf("--advertise-client-urls=%s", e.processState.URL),
		fmt.Sprintf("--listen-client-urls=%s", e.processState.URL),
		fmt.Sprintf("--data-dir=%s", e.processState.Dir),
	}

	return e.processState.Start(fmt.Sprintf("serving insecure client requests on %s", e.processState.URL.Hostname()))
}

// Stop stops this process gracefully, waits for its termination, and cleans up the data directory.
func (e *Etcd) Stop() error {
	return e.processState.Stop()
}
