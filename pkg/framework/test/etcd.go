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
	ProcessStarter simpleSessionStarter
	DataDirManager dataDirManager
	Config         *EtcdConfig
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

var etcdBinPathFinder = DefaultBinPathFinder

// NewEtcd returns a Etcd process configured with sane defaults
func NewEtcd() (*Etcd, error) {
	starter := func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
		return gexec.Start(command, out, err)
	}

	config, err := NewEtcdConfig()
	if err != nil {
		return nil, err
	}

	return &Etcd{
		Path:           etcdBinPathFinder("etcd"),
		ProcessStarter: starter,
		DataDirManager: NewTempDirManager(),
		Config:         config,
	}, nil
}

// NewEtcdWithBinaryAndConfig returns a Etcd process, using the handed in binary and config.
func NewEtcdWithBinaryAndConfig(pathToEtcd string, config *EtcdConfig) *Etcd {
	starter := func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
		return gexec.Start(command, out, err)
	}

	etcd := &Etcd{
		Path:           pathToEtcd,
		ProcessStarter: starter,
		DataDirManager: NewTempDirManager(),
		Config:         config,
	}

	return etcd
}

// GetURL returns the URL Etcd is listening on. Clients can use this to connect to Etcd.
func (e *Etcd) GetURL() string {
	return e.Config.ClientURL
}

// Start starts the etcd, waits for it to come up, and returns an error, if occoured.
func (e *Etcd) Start() error {
	if err := e.Config.Validate(); err != nil {
		return err
	}

	e.stdOut = gbytes.NewBuffer()
	e.stdErr = gbytes.NewBuffer()

	dataDir, err := e.DataDirManager.Create()
	if err != nil {
		return err
	}

	args := []string{
		"--debug",
		"--advertise-client-urls",
		e.Config.ClientURL,
		"--listen-client-urls",
		e.Config.ClientURL,
		"--listen-peer-urls",
		e.Config.PeerURL,
		"--data-dir",
		dataDir,
	}

	clientURL, err := url.Parse(e.Config.ClientURL)
	if err != nil {
		return err
	}

	detectedStart := e.stdErr.Detect(fmt.Sprintf(
		"serving insecure client requests on %s", clientURL.Host))
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
