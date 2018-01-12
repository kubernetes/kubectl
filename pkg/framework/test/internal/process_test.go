package internal_test

import (
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "k8s.io/kubectl/pkg/framework/test/internal"
)

var _ = Describe("Start", func() {
	var (
		command *exec.Cmd
		timeout time.Duration
	)
	BeforeEach(func() {
		command = getSimpleCommand()
		timeout = 200 * time.Millisecond
	})

	It("can start a process", func() {
		timeout = 5 * time.Second
		session, err := Start(command, "loop 3", timeout)

		Expect(err).NotTo(HaveOccurred())
		Consistently(session.ExitCode).Should(BeNumerically("==", -1))
	})

	Context("when process takes too long to start", func() {
		It("returns an timeout error", func() {
			session, err := Start(command, "loop 3000", timeout)

			Expect(err).To(MatchError(ContainSubstring("timeout")))
			Eventually(session.ExitCode, 10).Should(BeNumerically("==", 143))
		})
	})

	Context("when command cannot be started", func() {
		BeforeEach(func() {
			command = exec.Command("/notexistent")
		})
		It("propagates the error", func() {
			_, err := Start(command, "does not matter", timeout)

			Expect(os.IsNotExist(err)).To(BeTrue())
		})
	})
})

var _ = Describe("Stop", func() {
	var (
		session *gexec.Session
		timeout time.Duration
	)
	BeforeEach(func() {
		var err error
		command := getSimpleCommand()

		timeout = 10 * time.Second
		session, err = gexec.Start(command, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can stop a session", func() {
		err := Stop(session, timeout, "", false)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.ExitCode()).To(Equal(143))
	})

	Context("when command cannot be stopped", func() {
		BeforeEach(func() {
			session.Exited = make(chan struct{})
			timeout = 200 * time.Millisecond
		})
		It("runs into a timeout", func() {
			err := Stop(session, timeout, "", false)
			Expect(err).To(MatchError(ContainSubstring("timeout")))
		})
	})

	Context("when directory needs cleanup", func() {
		var (
			tmpDir string
		)
		BeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir("", "k8s_test_framework_tests_")
			Expect(err).NotTo(HaveOccurred())
			Expect(tmpDir).To(BeAnExistingFile())
		})
		It("removes the directory", func() {
			err := Stop(session, timeout, tmpDir, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(tmpDir).NotTo(BeAnExistingFile())
		})
	})
})

var _ = Describe("NewProcessState", func() {
	var (
		name         string
		path         string
		listenURL    *url.URL
		dir          string
		startTimeout time.Duration
		stopTimeout  time.Duration
	)
	BeforeEach(func() {
		name = "some name"
		path = "some path"
		listenURL = &url.URL{Host: "some host"}
		dir = "some dir"
		startTimeout = 1 * time.Second
		stopTimeout = 1 * time.Second
	})

	It("creates a new ProcessState struct", func() {
		procState, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)

		Expect(err).NotTo(HaveOccurred())
		Expect(procState.Path).To(Equal(path))
		Expect(procState.URL).To(Equal(listenURL))
		Expect(procState.Dir).To(Equal(dir))
		Expect(procState.StartTimeout).To(Equal(startTimeout))
		Expect(procState.StopTimeout).To(Equal(stopTimeout))
		Expect(procState.DirNeedsCleaning).To(BeFalse())
	})

	Context("When path is empty but symbolic name is set", func() {
		BeforeEach(func() {
			path = ""
		})
		It("defaults the path", func() {
			procState, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)

			Expect(err).NotTo(HaveOccurred())
			Expect(procState.Path).To(ContainSubstring("assets/bin/some name"))
		})
	})

	Context("When symbolic name is empty but path is set", func() {
		BeforeEach(func() {
			name = ""
		})
		It("defaults the path", func() {
			procState, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)

			Expect(err).NotTo(HaveOccurred())
			Expect(procState.Path).To(ContainSubstring("some path"))
		})
	})

	Context("When symbolic name and the path are empty", func() {
		BeforeEach(func() {
			name = ""
			path = ""
		})
		It("errors", func() {
			_, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)
			Expect(err).To(MatchError("Either a path or a symbolic name need to be set"))
		})
	})

	Context("When listen URL is not set", func() {
		BeforeEach(func() {
			listenURL = nil
		})
		It("defaults the URL", func() {
			procState, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)

			Expect(err).NotTo(HaveOccurred())
			Expect(procState.URL).NotTo(BeNil())
			Expect(procState.URL.Host).NotTo(BeEmpty())
		})
	})

	Context("When the directory is not set", func() {
		BeforeEach(func() {
			dir = ""
		})
		It("defaults to and creates a new directory and sets the directory cleanup flag", func() {
			procState, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)

			Expect(err).NotTo(HaveOccurred())
			Expect(procState.Dir).To(BeADirectory())
			Expect(procState.DirNeedsCleaning).To(BeTrue())

			Expect(os.RemoveAll(procState.Dir)).To(Succeed())
		})
	})

	Context("When timeouts are not set", func() {
		BeforeEach(func() {
			startTimeout = 0
			stopTimeout = 0
		})
		It("defaults both timeouts", func() {
			procState, err := NewProcessState(name, path, listenURL, dir, startTimeout, stopTimeout)

			Expect(err).NotTo(HaveOccurred())
			Expect(procState.StartTimeout).NotTo(BeZero())
			Expect(procState.StopTimeout).NotTo(BeZero())
		})
	})

})

func getSimpleCommand() *exec.Cmd {
	return exec.Command(
		"bash", "-c",
		`
			i=0
			while true
			do
				echo "loop $i" >&2
				let 'i += 1'
				sleep 0.2
			done
		`,
	)
}
