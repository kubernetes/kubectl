package test

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
	//TODO Stop should return an error
	Stop()
	URL() (string, error)
}

//go:generate counterfeiter . FixtureProcess

// NewFixtures will give you a Fixtures struct that's properly wired together.
func NewFixtures() (*Fixtures, error) {
	apiServer, err := NewAPIServer()
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
func (f *Fixtures) APIServerURL() (string, error) {
	return f.APIServer.URL()
}
