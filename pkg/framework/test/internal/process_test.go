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

var _ = Describe("Start method", func() {
	It("can start a process", func() {
		processState := &ProcessState{}
		processState.Path = "bash"
		processState.Args = simpleBashScript
		processState.StartTimeout = 10 * time.Second
		processState.StartMessage = "loop 5"

		err := processState.Start()
		Expect(err).NotTo(HaveOccurred())

		Consistently(processState.Session.ExitCode).Should(BeNumerically("==", -1))
	})

	Context("when process takes too long to start", func() {
		It("returns a timeout error", func() {
			processState := &ProcessState{}
			processState.Path = "bash"
			processState.Args = simpleBashScript
			processState.StartTimeout = 200 * time.Millisecond
			processState.StartMessage = "loop 5000"

			err := processState.Start()
			Expect(err).To(MatchError(ContainSubstring("timeout")))

			Eventually(processState.Session.ExitCode).Should(Equal(143))
		})
	})

	Context("when the command cannot be started", func() {
		It("propagates the error", func() {
			processState := &ProcessState{}
			processState.Path = "/nonexistent"

			err := processState.Start()

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

var _ = Describe("DoDefaulting", func() {
	Context("when all inputs are provided", func() {
		It("passes them through", func() {
			defaults, err := DoDefaulting(
				"some name",
				&url.URL{Host: "some.host.to.listen.on"},
				"/some/dir",
				"/some/path/to/some/bin",
				20*time.Hour,
				65537*time.Millisecond,
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(defaults.URL).To(Equal(url.URL{Host: "some.host.to.listen.on"}))
			Expect(defaults.Dir).To(Equal("/some/dir"))
			Expect(defaults.DirNeedsCleaning).To(BeFalse())
			Expect(defaults.Path).To(Equal("/some/path/to/some/bin"))
			Expect(defaults.StartTimeout).To(Equal(20 * time.Hour))
			Expect(defaults.StopTimeout).To(Equal(65537 * time.Millisecond))
		})
	})

	Context("when inputs are empty", func() {
		It("defaults them", func() {
			defaults, err := DoDefaulting(
				"some name",
				nil,
				"",
				"",
				0,
				0,
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(defaults.Dir).To(BeADirectory())
			Expect(os.RemoveAll(defaults.Dir)).To(Succeed())
			Expect(defaults.DirNeedsCleaning).To(BeTrue())

			Expect(defaults.URL).NotTo(BeZero())
			Expect(defaults.URL.Scheme).To(Equal("http"))
			Expect(defaults.URL.Hostname()).NotTo(BeEmpty())
			Expect(defaults.URL.Port()).NotTo(BeEmpty())

			Expect(defaults.Path).NotTo(BeEmpty())

			Expect(defaults.StartTimeout).NotTo(BeZero())
			Expect(defaults.StopTimeout).NotTo(BeZero())
		})
	})

	Context("when neither name nor path are provided", func() {
		It("returns an error", func() {
			_, err := DoDefaulting(
				"",
				nil,
				"",
				"",
				0,
				0,
			)
			Expect(err).To(MatchError("must have at least one of name or path"))
		})
	})
})

var simpleBashScript = []string{
	"-c",
	`
		i=0
		while true
		do
			echo "loop $i" >&2
			let 'i += 1'
			sleep 0.2
		done
	`,
}

func getSimpleCommand() *exec.Cmd {
	return exec.Command("bash", simpleBashScript...)
}
