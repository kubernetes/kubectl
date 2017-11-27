package test

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Etcd knows how to run an etcd server. Set it up with the path to a precompiled binary.
type Etcd struct {
	// The path to the etcd binary
	Path           string
	EtcdURL        string
	session        *gexec.Session
	stdOut         *gbytes.Buffer
	stdErr         *gbytes.Buffer
	dataDirManager DataDirManager
}

// DataDirManager knows how to create and destroy Etcd's data directory.
type DataDirManager interface {
	Create() (string, error)
	Destroy() error
}

// Start starts the etcd, and returns a gexec.Session. To stop it again, call Terminate and Wait on that session.
func (e *Etcd) Start() error {
	e.dataDirManager = NewTempDirManager()
	e.stdOut = gbytes.NewBuffer()
	e.stdErr = gbytes.NewBuffer()

	dataDir, err := e.dataDirManager.Create()
	if err != nil {
		return err
	}

	args := []string{
		"--debug",
		"--advertise-client-urls",
		e.EtcdURL,
		"--listen-client-urls",
		e.EtcdURL,
		"--data-dir",
		dataDir,
	}

	detectedStart := e.stdErr.Detect("serving insecure client requests on 127.0.0.1:2379")
	timedOut := time.After(20 * time.Second)

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

// Stop stops this process gracefully.
func (e *Etcd) Stop() {
	if e.session != nil {
		e.session.Terminate().Wait(20 * time.Second)
		err := e.dataDirManager.Destroy()
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
