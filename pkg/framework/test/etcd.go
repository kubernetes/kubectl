package test

import (
	"fmt"
	"os/exec"
	"time"

	"net/url"

	"github.com/onsi/gomega/gexec"
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
	session      *gexec.Session

	commonStuff internal.CommonStuff
}

// Start starts the etcd, waits for it to come up, and returns an error, if occoured.
func (e *Etcd) Start() error {
	err := e.ensureInitialized()
	if err != nil {
		return err
	}

	args := []string{
		"--debug",
		"--listen-peer-urls=http://localhost:0",
		fmt.Sprintf("--advertise-client-urls=%s", e.commonStuff.URL),
		fmt.Sprintf("--listen-client-urls=%s", e.commonStuff.URL),
		fmt.Sprintf("--data-dir=%s", e.commonStuff.Dir),
	}

	e.session, err = internal.Start(
		exec.Command(e.commonStuff.Path, args...),
		fmt.Sprintf("serving insecure client requests on %s", e.commonStuff.URL.Hostname()),
		e.commonStuff.StartTimeout,
	)

	return err
}

func (e *Etcd) ensureInitialized() error {
	var err error
	e.commonStuff, err = internal.NewCommonStuff(
		"etcd",
		e.Path,
		e.URL,
		e.DataDir,
		e.StartTimeout, e.StopTimeout,
	)
	return err
}

// Stop stops this process gracefully, waits for its termination, and cleans up the data directory.
func (e *Etcd) Stop() error {
	return internal.Stop(
		e.session,
		e.commonStuff.StopTimeout,
		e.commonStuff.Dir,
		e.commonStuff.DirNeedsCleaning,
	)
}
