package test

import (
	"fmt"
	"net"
)

// Fixtures is a struct that knows how to start all your test fixtures.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in future.
type Fixtures struct {
	Etcd      FixtureProcess
	APIServer FixtureProcess
	Config    FixturesConfig
}

// FixturesConfig is a datastructure that exposes configuration that should be used by clients to talk
// to the fixture processes.
type FixturesConfig struct {
	APIServerURL string
}

// FixtureProcess knows how to start and stop a Fixture processes.
// This interface is potentially going to be expanded to e.g. allow access to the processes StdOut/StdErr
// and other internals.
type FixtureProcess interface {
	Start() error
	Stop()
}

//go:generate counterfeiter . FixtureProcess

// NewFixtures will give you a Fixtures struct that's properly wired together.
func NewFixtures(pathToEtcd, pathToAPIServer string) (*Fixtures, error) {
	urls := map[string]string{
		"etcdClients":      "",
		"etcdPeers":        "",
		"apiServerClients": "",
	}
	host := "127.0.0.1"

	for name := range urls {
		port, err := getFreePort(host)
		if err != nil {
			return nil, err
		}
		urls[name] = fmt.Sprintf("http://%s:%d", host, port)
	}

	fixtures := &Fixtures{
		Etcd:      NewEtcd(pathToEtcd, urls["etcdClients"], urls["etcdPeers"]),
		APIServer: NewAPIServer(pathToAPIServer, urls["apiServerClients"], urls["etcdClients"]),
	}

	fixtures.Config = FixturesConfig{
		APIServerURL: urls["apiServerClients"],
	}

	return fixtures, nil
}

// Start will start all your fixtures. To stop them, call Stop().
func (f *Fixtures) Start() error {
	started := make(chan error)
	starter := func(process FixtureProcess) {
		started <- process.Start()
	}
	processes := []FixtureProcess{
		f.Etcd,
		f.APIServer,
	}

	for _, process := range processes {
		go starter(process)
	}

	for pendingProcesses := len(processes); pendingProcesses > 0; pendingProcesses-- {
		if err := <-started; err != nil {
			return err
		}
	}

	return nil
}

// Stop will stop all your fixtures, and clean up their data.
func (f *Fixtures) Stop() error {
	f.APIServer.Stop()
	f.Etcd.Stop()
	return nil
}

func getFreePort(host string) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", host))
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
