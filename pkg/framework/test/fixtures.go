package test

// Fixtures is a struct that knows how to start all your test fixtures.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in future.
type Fixtures struct {
	Etcd      FixtureProcess
	APIServer FixtureProcess
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
func NewFixtures(pathToEtcd, pathToAPIServer string) *Fixtures {
	etcdURL := "http://127.0.0.1:2379"
	return &Fixtures{
		Etcd: &Etcd{
			Path:    pathToEtcd,
			EtcdURL: etcdURL,
		},
		APIServer: &APIServer{
			Path:    pathToAPIServer,
			EtcdURL: etcdURL,
		},
	}
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
