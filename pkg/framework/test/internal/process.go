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

//type ProcessState2 struct {
//	ProcessInput
//	Args        []string
//	StartString string
//}
//
//func NewProcessState(input ProcessInput, args []string, startthing string) ProcessState {
//	return ProcessState2{input, args, startthing}
//}

func (ps *ProcessState) Start() (err error) {
	ps.Session, err = Start(
		ps.Path,
		ps.Args,
		ps.StartMessage,
		ps.StartTimeout,
	)
	return
}

func (ps *ProcessState) Stop() error {
	return Stop(
		ps.Session,
		ps.StopTimeout,
		ps.Dir,
		ps.DirNeedsCleaning,
	)
}

func Start(path string, args []string, startMessage string, startTimeout time.Duration) (*gexec.Session, error) {
	command := exec.Command(path, args...)

	stdErr := gbytes.NewBuffer()
	detectedStart := stdErr.Detect(startMessage)
	timedOut := time.After(startTimeout)

	session, err := gexec.Start(command, nil, stdErr)
	if err != nil {
		return session, err
	}

	select {
	case <-detectedStart:
		return session, nil
	case <-timedOut:
		session.Terminate()
		return session, fmt.Errorf("timeout waiting for process to start serving")
	}

}

func Stop(session *gexec.Session, stopTimeout time.Duration, dirToClean string, dirNeedsCleaning bool) error {
	if session == nil {
		return nil
	}

	detectedStop := session.Terminate().Exited
	timedOut := time.After(stopTimeout)

	select {
	case <-detectedStop:
		break
	case <-timedOut:
		return fmt.Errorf("timeout waiting for process to stop")
	}

	if dirNeedsCleaning {
		return os.RemoveAll(dirToClean)
	}

	return nil
}
