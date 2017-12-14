package test

import (
	"fmt"
	"net"
)

// Fixtures is a struct that knows how to start all your test fixtures.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in future.
type Fixtures struct {
	APIServer FixtureProcess
}

// FixtureProcess knows how to start and stop a Fixture processes.
// This interface is potentially going to be expanded to e.g. allow access to the processes StdOut/StdErr
// and other internals.
type FixtureProcess interface {
	Start() error
	Stop()
	URL() string
}

//go:generate counterfeiter . FixtureProcess

// NewFixtures will give you a Fixtures struct that's properly wired together.
func NewFixtures() (*Fixtures, error) {
	apiServerConfig := &APIServerConfig{}

	if url, urlErr := getHTTPListenURL(); urlErr == nil {
		apiServerConfig.APIServerURL = url
	} else {
		return nil, urlErr
	}

	apiServer, err := NewAPIServer(apiServerConfig)
	if err != nil {
		return nil, err
	}

	fixtures := &Fixtures{
		APIServer: apiServer,
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
		f.APIServer,
	}

	for _, process := range processes {
		go starter(process)
	}

	for range processes {
		if err := <-started; err != nil {
			return err
		}
	}

	return nil
}

// Stop will stop all your fixtures, and clean up their data.
func (f *Fixtures) Stop() error {
	f.APIServer.Stop()
	return nil
}

// APIServerURL returns the URL to the APIServer. Clients can use this URL to connect to the APIServer.
func (f *Fixtures) APIServerURL() string {
	return f.APIServer.URL()
}

func getHTTPListenURL() (url string, err error) {
	host := "127.0.0.1"
	port, err := getFreePort(host)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d", host, port), nil
}

func getFreePort(host string) (port int, err error) {
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", host))
	if err != nil {
		return
	}
	var l *net.TCPListener
	l, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer func() {
		err = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port, nil
}
