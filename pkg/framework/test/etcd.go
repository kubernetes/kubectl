package test

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Etcd knows how to run an etcd server. Set it up with the path to a precompiled binary.
type Etcd struct {
	AddressManager AddressManager
	PathFinder     BinPathFinder
	ProcessStarter simpleSessionStarter
	DataDirManager dataDirManager
	session        SimpleSession
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
}

type dataDirManager interface {
	Create() (string, error)
	Destroy() error
}

//go:generate counterfeiter . dataDirManager

// SimpleSession describes a CLI session. You can get output, and you can kill it. It is implemented by *gexec.Session.
type SimpleSession interface {
	Buffer() *gbytes.Buffer
	ExitCode() int
	Wait(timeout ...interface{}) *gexec.Session
	Terminate() *gexec.Session
}

//go:generate counterfeiter . SimpleSession

type simpleSessionStarter func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error)

// URL returns the URL Etcd is listening on. Clients can use this to connect to Etcd.
func (e *Etcd) URL() (string, error) {
	port, err := e.AddressManager.Port()
	if err != nil {
		return "", err
	}
	host, err := e.AddressManager.Host()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d", host, port), nil
}

// Start starts the etcd, waits for it to come up, and returns an error, if occoured.
func (e *Etcd) Start() error {
	e.ensureInitialized()

	port, host, err := e.AddressManager.Initialize("localhost")
	if err != nil {
		return err
	}

	dataDir, err := e.DataDirManager.Create()
	if err != nil {
		return err
	}

	clientURL := fmt.Sprintf("http://%s:%d", host, port)
	args := []string{
		"--debug",
		"--listen-peer-urls=http://localhost:0",
		fmt.Sprintf("--advertise-client-urls=%s", clientURL),
		fmt.Sprintf("--listen-client-urls=%s", clientURL),
		fmt.Sprintf("--data-dir=%s", dataDir),
	}

	detectedStart := e.stdErr.Detect(fmt.Sprintf(
		"serving insecure client requests on %s", host))
	timedOut := time.After(20 * time.Second)

	command := exec.Command(e.PathFinder("etcd"), args...)
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

func (e *Etcd) ensureInitialized() {
	if e.PathFinder == nil {
		e.PathFinder = DefaultBinPathFinder
	}

	if e.AddressManager == nil {
		e.AddressManager = &DefaultAddressManager{}
	}
	if e.ProcessStarter == nil {
		e.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
			return gexec.Start(command, out, err)
		}
	}
	if e.DataDirManager == nil {
		e.DataDirManager = NewTempDirManager()
	}

	e.stdOut = gbytes.NewBuffer()
	e.stdErr = gbytes.NewBuffer()
}

// Stop stops this process gracefully, waits for its termination, and cleans up the data directory.
func (e *Etcd) Stop() error {
	if e.session == nil {
		return nil
	}

	e.session.Terminate()
	// TODO have a better way to handle the timeout of Stop()
	e.session.Wait(20 * time.Second)

	err := e.DataDirManager.Destroy()

	return err
}

// ExitCode returns the exit code of the process, if it has exited. If it hasn't exited yet, ExitCode returns -1.
func (e *Etcd) ExitCode() int {
	return e.session.ExitCode()
}

// Buffer implements the gbytes.BufferProvider interface and returns the stdout of the process
func (e *Etcd) Buffer() *gbytes.Buffer {
	return e.session.Buffer()
}
