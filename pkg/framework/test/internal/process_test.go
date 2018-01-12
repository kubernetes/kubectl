package internal_test

import (
	. "k8s.io/kubectl/pkg/framework/test/internal"

	"os/exec"

	"time"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
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
