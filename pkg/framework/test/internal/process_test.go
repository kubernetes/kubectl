package internal_test

import (
	. "k8s.io/kubectl/pkg/framework/test/internal"

	"os/exec"

	"time"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
