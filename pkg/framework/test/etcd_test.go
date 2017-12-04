package test_test

import (
	"io"
	"os/exec"

	. "k8s.io/kubectl/pkg/framework/test"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"k8s.io/kubectl/pkg/framework/test/testfakes"
)

var _ = Describe("Etcd", func() {
	var (
		fakeSession *testfakes.FakeSimpleSession
		etcd        *Etcd
	)

	BeforeEach(func() {
		fakeSession = &testfakes.FakeSimpleSession{}
		etcd = &Etcd{
			Path:    "",
			EtcdURL: "our etcd url",
		}
	})

	Context("when given a path to a binary that runs for a long time", func() {
		It("can start and stop that binary", func() {
			sessionBuffer := gbytes.NewBuffer()
			fmt.Fprintf(sessionBuffer, "Everything is dandy")
			fakeSession.BufferReturns(sessionBuffer)

			fakeSession.ExitCodeReturnsOnCall(0, -1)
			fakeSession.ExitCodeReturnsOnCall(1, 143)

			etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				fmt.Fprint(err, "serving insecure client requests on 127.0.0.1:2379")
				return fakeSession, nil
			}

			By("Starting the Etcd Server")
			err := etcd.Start()
			Expect(err).NotTo(HaveOccurred())

			Eventually(etcd).Should(gbytes.Say("Everything is dandy"))
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(0))
			Expect(etcd).NotTo(gexec.Exit())
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(1))

			By("Stopping the Etcd Server")
			etcd.Stop()
			Expect(etcd).To(gexec.Exit(143))
			Expect(fakeSession.TerminateCallCount()).To(Equal(1))
			Expect(fakeSession.WaitCallCount()).To(Equal(1))
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(2))
		})
	})

	Context("when  the starter returns an error", func() {
		It("passes the error to the caller", func() {
			etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				return nil, fmt.Errorf("Some error in the starter.")
			}

			err := etcd.Start()
			Expect(err).To(MatchError(ContainSubstring("Some error in the starter.")))
		})
	})

	Context("when we try to stop a server that hasn't been started", func() {
		It("is a noop and does not call exit on the session", func() {
			etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				return fakeSession, nil
			}
			etcd.Stop()
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(0))
		})
	})
})
