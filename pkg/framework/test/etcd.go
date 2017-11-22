package test

import (
	"os/exec"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
)

// Etcd knows how to run an etcd server. Set it up with the path to a precompiled binary.
type Etcd struct {
	// The path to the etcd binary
	Path string
}

// Start starts the etcd, and returns a gexec.Session. To stop it again, call Terminate and Wait on that session.
func (s Etcd) Start(etcdURL string, datadir string) (*gexec.Session, error) {
	args := []string{
		"--advertise-client-urls",
		etcdURL,
		"--data-dir",
		datadir,
		"--listen-client-urls",
		etcdURL,
		"--debug",
	}

	command := exec.Command(s.Path, args...)
	return gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
}
