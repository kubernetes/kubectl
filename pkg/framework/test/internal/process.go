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
	URL              *url.URL
	Dir              string
	DirNeedsCleaning bool
	Path             string
	StopTimeout      time.Duration
	StartTimeout     time.Duration
}

func Start(command *exec.Cmd, startMessage string, startTimeout time.Duration) (*gexec.Session, error) {
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

func NewProcessState(
	symbolicName string,
	path string,
	listenURL *url.URL,
	dir string,
	startTimeout time.Duration,
	stopTimeout time.Duration,
) (ProcessState, error) {
	if path == "" && symbolicName == "" {
		return ProcessState{}, fmt.Errorf("Either a path or a symbolic name need to be set")
	}

	state := ProcessState{
		Path:             path,
		URL:              listenURL,
		Dir:              dir,
		DirNeedsCleaning: false,
		StartTimeout:     startTimeout,
		StopTimeout:      stopTimeout,
	}

	if path == "" {
		state.Path = BinPathFinder(symbolicName)
	}

	if listenURL == nil {
		am := &AddressManager{}
		port, host, err := am.Initialize()
		if err != nil {
			return ProcessState{}, err
		}
		state.URL = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		}
	}

	if dir == "" {
		newDir, err := ioutil.TempDir("", "k8s_test_framework_")
		if err != nil {
			return ProcessState{}, err
		}
		state.Dir = newDir
		state.DirNeedsCleaning = true
	}

	if stopTimeout == 0 {
		state.StopTimeout = 20 * time.Second
	}

	if startTimeout == 0 {
		state.StartTimeout = 20 * time.Second
	}

	return state, nil
}
