package test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"net/url"

	"os"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"k8s.io/kubectl/pkg/framework/test/internal"
)

// Etcd knows how to run an etcd server.
//
// The documentation and examples for the Etcd's properties can be found in
// in the documentation for the `APIServer`, as both implement a `ControlPaneProcess`.
type Etcd struct {
	Address       *url.URL
	Path          string
	DataDir       string
	actualDataDir string
	StopTimeout   time.Duration
	StartTimeout  time.Duration
	session       *gexec.Session
	stdOut        *gbytes.Buffer
	stdErr        *gbytes.Buffer
}

// URL returns the URL Etcd is listening on. Clients can use this to connect to Etcd.
func (e *Etcd) URL() (string, error) {
	if e.Address == nil {
		return "", fmt.Errorf("Etcd's Address not initialized or configured")
	}
	return e.Address.String(), nil
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
		fmt.Sprintf("--advertise-client-urls=%s", e.Address),
		fmt.Sprintf("--listen-client-urls=%s", e.Address),
		fmt.Sprintf("--data-dir=%s", e.actualDataDir),
	}

	detectedStart := e.stdErr.Detect(fmt.Sprintf(
		"serving insecure client requests on %s", e.Address.Hostname()))
	timedOut := time.After(e.StartTimeout)

	command := exec.Command(e.Path, args...)
	e.session, err = gexec.Start(command, e.stdOut, e.stdErr)
	if err != nil {
		return err
	}

	select {
	case <-detectedStart:
		return nil
	case <-timedOut:
		return fmt.Errorf("timeout waiting for etcd to start serving")
	}
}

func (e *Etcd) ensureInitialized() error {
	if e.Path == "" {
		e.Path = internal.BinPathFinder("etcd")
	}
	if e.Address == nil {
		am := &internal.AddressManager{}
		port, host, err := am.Initialize()
		if err != nil {
			return err
		}

		e.Address = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		}
	}
	if e.DataDir == "" {
		dataDir, err := ioutil.TempDir("", "k8s_test_framework_")
		if err != nil {
			return err
		}
		e.actualDataDir = dataDir
	} else {
		e.actualDataDir = e.DataDir
	}

	if e.StopTimeout == 0 {
		e.StopTimeout = 20 * time.Second
	}
	if e.StartTimeout == 0 {
		e.StartTimeout = 20 * time.Second
	}

	e.stdOut = gbytes.NewBuffer()
	e.stdErr = gbytes.NewBuffer()

	return nil
}

// Stop stops this process gracefully, waits for its termination, and cleans up the data directory.
func (e *Etcd) Stop() error {
	if e.session == nil {
		return nil
	}

	session := e.session.Terminate()
	detectedStop := session.Exited
	timedOut := time.After(e.StopTimeout)

	select {
	case <-detectedStop:
		break
	case <-timedOut:
		return fmt.Errorf("timeout waiting for etcd to stop")
	}

	if e.DataDir == "" {
		return os.RemoveAll(e.actualDataDir)
	}

	return nil
}
