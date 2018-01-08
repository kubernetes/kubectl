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

// Etcd knows how to run an etcd server.
//
// The documentation and examples for the Etcd's properties can be found in
// in the documentation for the `APIServer`, as both implement a `ControlPaneProcess`.
type Etcd struct {
	Address        *url.URL
	Path           string
	ProcessStarter SimpleSessionStarter
	DataDir        *CleanableDirectory
	StopTimeout    time.Duration
	StartTimeout   time.Duration
	session        SimpleSession
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
}

// SimpleSession describes a CLI session. You can get output, the exit code, and you can terminate it.
//
// It is implemented by *gexec.Session.
type SimpleSession interface {
	Buffer() *gbytes.Buffer
	ExitCode() int
	Terminate() *gexec.Session
}

//go:generate counterfeiter . SimpleSession

// SimpleSessionStarter knows how to start a exec.Cmd with a writer for both StdOut & StdErr and returning it wrapped
// in a `SimpleSession`.
type SimpleSessionStarter func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error)

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
		fmt.Sprintf("--data-dir=%s", e.DataDir.Path),
	}

	detectedStart := e.stdErr.Detect(fmt.Sprintf(
		"serving insecure client requests on %s", e.Address.Hostname()))
	timedOut := time.After(e.StartTimeout)

	command := exec.Command(e.Path, args...)
	e.session, err = e.ProcessStarter(command, e.stdOut, e.stdErr)
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
	if e.ProcessStarter == nil {
		e.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
			return gexec.Start(command, out, err)
		}
	}
	if e.DataDir == nil {
		dataDir, err := newDirectory()
		if err != nil {
			return err
		}
		e.DataDir = dataDir
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

	if e.DataDir.Cleanup == nil {
		return nil
	}
	return e.DataDir.Cleanup()
}

// Buffer implements the gbytes.BufferProvider interface and returns the stdout of the process
func (e *Etcd) Buffer() *gbytes.Buffer {
	return e.session.Buffer()
}
