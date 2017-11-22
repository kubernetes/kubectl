package test

import (
	"os/exec"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Etcd knows how to run an etcd server. Set it up with the path to a precompiled binary.
type Etcd struct {
	// The path to the etcd binary
	Path    string
	session *gexec.Session
}

// Start starts the etcd, and returns a gexec.Session. To stop it again, call Terminate and Wait on that session.
func (e *Etcd) Start(etcdURL string, datadir string) error {
	args := []string{
		"--advertise-client-urls",
		etcdURL,
		"--data-dir",
		datadir,
		"--listen-client-urls",
		etcdURL,
		"--debug",
	}

	command := exec.Command(e.Path, args...)
	var err error
	e.session, err = gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	return err
}

// Stop stops this process gracefully.
func (e *Etcd) Stop() {
	if e.session != nil {
		e.session.Terminate()
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
