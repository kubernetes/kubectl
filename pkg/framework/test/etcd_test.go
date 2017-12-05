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
		fakeSession        *testfakes.FakeSimpleSession
		fakeDataDirManager *testfakes.FakeDataDirManager
		etcd               *Etcd
		etcdConfig         map[string]string
	)

	BeforeEach(func() {
		fakeSession = &testfakes.FakeSimpleSession{}
		fakeDataDirManager = &testfakes.FakeDataDirManager{}

		etcd = &Etcd{
			Path:           "",
			DataDirManager: fakeDataDirManager,
		}

		etcdConfig = map[string]string{
			"clientURL": "http://this.is.etcd.listening.for.clients:1234",
			"peerURL":   "http://this.is.etcd.listening.for.peers:1235",
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
				fmt.Fprint(err, "serving insecure client requests on this.is.etcd.listening.for.clients:1234")
				return fakeSession, nil
			}

			By("Starting the Etcd Server")
			err := etcd.Start(etcdConfig)
			Expect(err).NotTo(HaveOccurred())

			Eventually(etcd).Should(gbytes.Say("Everything is dandy"))
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(0))
			Expect(etcd).NotTo(gexec.Exit())
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(1))
			Expect(fakeDataDirManager.CreateCallCount()).To(Equal(1))

			By("Stopping the Etcd Server")
			etcd.Stop()
			Expect(etcd).To(gexec.Exit(143))
			Expect(fakeSession.TerminateCallCount()).To(Equal(1))
			Expect(fakeSession.WaitCallCount()).To(Equal(1))
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(2))
			Expect(fakeDataDirManager.DestroyCallCount()).To(Equal(1))
		})
	})

	Context("when the data directory cannot be created", func() {
		It("propagates the error, and does not start the process", func() {
			fakeDataDirManager.CreateReturnsOnCall(0, "", fmt.Errorf("Error on directory creation."))

			processStarterCounter := 0
			etcd.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				processStarterCounter += 1
				return fakeSession, nil
			}

			err := etcd.Start(etcdConfig)
			Expect(err).To(MatchError(ContainSubstring("Error on directory creation.")))
			Expect(processStarterCounter).To(Equal(0))
		})
	})

	Context("when  the starter returns an error", func() {
		It("propagates the error", func() {
			etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				return nil, fmt.Errorf("Some error in the starter.")
			}

			err := etcd.Start(etcdConfig)
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
