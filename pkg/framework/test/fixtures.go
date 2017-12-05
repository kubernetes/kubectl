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
	URLGetter listenURLGetter
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
	Start(config map[string]string) error
	Stop()
}

//go:generate counterfeiter . FixtureProcess

// NewFixtures will give you a Fixtures struct that's properly wired together.
func NewFixtures(pathToEtcd, pathToAPIServer string) *Fixtures {
	fixtures := &Fixtures{
		Etcd:      NewEtcd(pathToEtcd),
		APIServer: NewAPIServer(pathToAPIServer),
		URLGetter: getHTTPListenURL,
	}

	return fixtures
}

// Start will start all your fixtures. To stop them, call Stop().
func (f *Fixtures) Start() error {
	type configs map[string]string

	etcdClientURL, err := f.URLGetter()
	if err != nil {
		return err
	}
	etcdPeerURL, err := f.URLGetter()
	if err != nil {
		return err
	}
	apiServerURL, err := f.URLGetter()
	if err != nil {
		return err
	}

	etcdConf := configs{
		"peerURL":   etcdPeerURL,
		"clientURL": etcdClientURL,
	}
	apiServerConf := configs{
		"etcdURL":      etcdClientURL,
		"apiServerURL": apiServerURL,
	}

	started := make(chan error)
	starter := func(process FixtureProcess, conf configs) {
		started <- process.Start(conf)
	}
	processes := map[FixtureProcess]configs{
		f.Etcd:      etcdConf,
		f.APIServer: apiServerConf,
	}

	for process, config := range processes {
		go starter(process, config)
	}

	for range processes {
		if err := <-started; err != nil {
			return err
		}
	}

	f.Config = FixturesConfig{
		APIServerURL: apiServerURL,
	}

	return nil
}

// Stop will stop all your fixtures, and clean up their data.
func (f *Fixtures) Stop() error {
	f.APIServer.Stop()
	f.Etcd.Stop()
	return nil
}

type listenURLGetter func() (url string, err error)

//go:generate counterfeiter . listenURLGetter

func getHTTPListenURL() (url string, err error) {
	host := "127.0.0.1"
	port, err := getFreePort(host)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d", host, port), nil
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
