package test

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"net/url"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Etcd knows how to run an etcd server. Set it up with the path to a precompiled binary.
type Etcd struct {
	Path           string
	EtcdURL        string
	EtcdPeerURL    string
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

// NewEtcd constructs an Etcd Fixture Process
func NewEtcd(pathToEtcd, etcdURL, etcdPeerURL string) *Etcd {
	starter := func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
		return gexec.Start(command, out, err)
	}

	etcd := &Etcd{
		Path:           pathToEtcd,
		EtcdURL:        etcdURL,
		EtcdPeerURL:    etcdPeerURL,
		ProcessStarter: starter,
		DataDirManager: NewTempDirManager(),
	}

	return etcd
}

// Start starts the etcd, waits for it to come up, and returns an error, if occoured.
func (e *Etcd) Start() error {
	e.stdOut = gbytes.NewBuffer()
	e.stdErr = gbytes.NewBuffer()

	dataDir, err := e.DataDirManager.Create()
	if err != nil {
		return err
	}

	args := []string{
		"--debug",
		"--advertise-client-urls",
		e.EtcdURL,
		"--listen-client-urls",
		e.EtcdURL,
		"--listen-peer-urls",
		e.EtcdPeerURL,
		"--data-dir",
		dataDir,
	}

	url, err := url.Parse(e.EtcdURL)
	if err != nil {
		return err
	}

	detectedStart := e.stdErr.Detect(fmt.Sprintf(
		"serving insecure client requests on %s", url.Host))
	timedOut := time.After(20 * time.Second)

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

// Stop stops this process gracefully, waits for its termination, and cleans up the data directory.
func (e *Etcd) Stop() {
	if e.session != nil {
		e.session.Terminate()
		e.session.Wait(20 * time.Second)
		err := e.DataDirManager.Destroy()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
}

// ExitCode returns the exit code of the process, if it has exited. If it hasn't exited yet, ExitCode returns -1.
func (e *Etcd) ExitCode() int {
	return e.session.ExitCode()
}

// Buffer implements the gbytes.BufferProvider interface and returns the stdout of the process
func (e *Etcd) Buffer() *gbytes.Buffer {
	return e.session.Buffer()
}
