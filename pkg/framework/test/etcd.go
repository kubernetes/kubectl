package test

import (
	"io/ioutil"
	"os"
	"os/exec"

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
	tempDirManager TempDirManager
}

// Start starts the etcd, and returns a gexec.Session. To stop it again, call Terminate and Wait on that session.
func (e *Etcd) Start() error {
	e.tempDirManager = &tempDirManager{}

	dataDir := e.tempDirManager.Create()

	args := []string{
		"--debug",
		"--advertise-client-urls",
		e.EtcdURL,
		"--listen-client-urls",
		e.EtcdURL,
		"--data-dir",
		dataDir,
	}

	command := exec.Command(e.Path, args...)
	var err error
	e.session, err = gexec.Start(command, e.stdOut, e.stdErr)
	return err
}

// Stop stops this process gracefully.
func (e *Etcd) Stop() {
	if e.session != nil {
		e.session.Terminate().Wait()
		err := e.tempDirManager.Destroy()
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

//------

// TempDirManager knows how to create and destroy temporary directories.
type TempDirManager interface {
	Create() string
	Destroy() error
}

//go:generate counterfeiter . TempDirManager

type tempDirManager struct {
	dir string
}

func (t *tempDirManager) Create() string {
	var err error
	t.dir, err = ioutil.TempDir("", "kube-test-framework")
	gomega.ExpectWithOffset(2, err).NotTo(gomega.HaveOccurred(),
		"expected to be able to create a temporary directory in the kube test framework")
	return t.dir
}

func (t *tempDirManager) Destroy() error {
	if t.dir != "" {
		return os.RemoveAll(t.dir)
	}
	return nil
}
