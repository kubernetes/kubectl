package internal

import (
	"fmt"
	"net/url"
	"os/exec"
	"time"

	"os"

	"io/ioutil"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type ProcessState struct {
	DefaultedProcessInput
	Session      *gexec.Session // TODO private?
	StartMessage string
	Args         []string
}

// TODO explore ProcessInputs, Defaulter, ProcessState, ...
type DefaultedProcessInput struct {
	URL              url.URL
	Dir              string
	DirNeedsCleaning bool
	Path             string
	StopTimeout      time.Duration
	StartTimeout     time.Duration
}

func DoDefaulting(
	name string,
	listenUrl *url.URL,
	dir string,
	path string,
	startTimeout time.Duration,
	stopTimeout time.Duration,
) (DefaultedProcessInput, error) {
	defaults := DefaultedProcessInput{
		Dir:          dir,
		Path:         path,
		StartTimeout: startTimeout,
		StopTimeout:  stopTimeout,
	}

	if listenUrl == nil {
		am := &AddressManager{}
		port, host, err := am.Initialize()
		if err != nil {
			return DefaultedProcessInput{}, err
		}
		defaults.URL = url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		}
	} else {
		defaults.URL = *listenUrl
	}

	if dir == "" {
		newDir, err := ioutil.TempDir("", "k8s_test_framework_")
		if err != nil {
			return DefaultedProcessInput{}, err
		}
		defaults.Dir = newDir
		defaults.DirNeedsCleaning = true
	}

	if path == "" {
		if name == "" {
			return DefaultedProcessInput{}, fmt.Errorf("must have at least one of name or path")
		}
		defaults.Path = BinPathFinder(name)
	}

	if startTimeout == 0 {
		defaults.StartTimeout = 20 * time.Second
	}

	if stopTimeout == 0 {
		defaults.StopTimeout = 20 * time.Second
	}

	return defaults, nil
}

func (ps *ProcessState) Start() (err error) {
	command := exec.Command(ps.Path, ps.Args...)

	stdErr := gbytes.NewBuffer()
	detectedStart := stdErr.Detect(ps.StartMessage)
	timedOut := time.After(ps.StartTimeout)

	ps.Session, err = gexec.Start(command, nil, stdErr)
	if err != nil {
		return err
	}

	select {
	case <-detectedStart:
		return nil
	case <-timedOut:
		ps.Session.Terminate()
		return fmt.Errorf("timeout waiting for process to start serving")
	}
}

func (ps *ProcessState) Stop() error {
	if ps.Session == nil {
		return nil
	}

	detectedStop := ps.Session.Terminate().Exited
	timedOut := time.After(ps.StopTimeout)

	select {
	case <-detectedStop:
		break
	case <-timedOut:
		return fmt.Errorf("timeout waiting for process to stop")
	}

	if ps.DirNeedsCleaning {
		return os.RemoveAll(ps.Dir)
	}

	return nil
}
