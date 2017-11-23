package test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/onsi/gomega"
)

// Fixtures is a struct that knows how to start all your test fixtures.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in future.
type Fixtures struct {
	Etcd           EtcdStartStopper
	APIServer      APIServerStartStopper
	TempDirManager TempDirManager
}

// EtcdStartStopper knows how to start an Etcd. One good implementation is Etcd.
type EtcdStartStopper interface {
	Start(etcdURL, datadir string) error
	Stop()
}

//go:generate counterfeiter . EtcdStartStopper

// APIServerStartStopper knows how to start an APIServer. One good implementation is APIServer.
type APIServerStartStopper interface {
	Start(etcdURL string) error
	Stop()
}

//go:generate counterfeiter . APIServerStartStopper

// TempDirManager knows how to create and destroy temporary directories.
type TempDirManager interface {
	Create() string
	Destroy() error
}

//go:generate counterfeiter . TempDirManager

// NewFixtures will give you a Fixtures struct that's properly wired together.
func NewFixtures(pathToEtcd, pathToAPIServer string) *Fixtures {
	return &Fixtures{
		Etcd:           &Etcd{Path: pathToEtcd},
		APIServer:      &APIServer{Path: pathToAPIServer},
		TempDirManager: &tempDirManager{},
	}
}

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

// Start will start all your fixtures. To stop them, call Stop().
func (f *Fixtures) Start() error {
	tmpDir := f.TempDirManager.Create()
	if err := f.Etcd.Start("http://127.0.0.1:2379", tmpDir); err != nil {
		return fmt.Errorf("Error starting etcd: %s", err)
	}
	if err := f.APIServer.Start("http://127.0.0.1:2379"); err != nil {
		return fmt.Errorf("Error starting apiserver: %s", err)
	}
	return nil
}

// Stop will stop all your fixtures, and clean up their data.
func (f *Fixtures) Stop() error {
	f.APIServer.Stop()
	f.Etcd.Stop()
	return f.TempDirManager.Destroy()
}
