package test_test

import (
	"fmt"
	"io"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "k8s.io/kubectl/pkg/framework/test"
	"k8s.io/kubectl/pkg/framework/test/testfakes"
)

var _ = Describe("Etcd", func() {
	var (
		fakeSession        *testfakes.FakeSimpleSession
		fakeDataDirManager *testfakes.FakeDataDirManager
		fakePathFinder     *testfakes.FakeBinPathFinder
		etcd               *Etcd
		etcdConfig         *EtcdConfig
	)

	BeforeEach(func() {
		fakeSession = &testfakes.FakeSimpleSession{}
		fakeDataDirManager = &testfakes.FakeDataDirManager{}
		fakePathFinder = &testfakes.FakeBinPathFinder{}

		etcdConfig = &EtcdConfig{
			ClientURL: "http://this.is.etcd.listening.for.clients:1234",
			PeerURL:   "http://this.is.etcd.listening.for.peers:1235",
		}

		etcd = &Etcd{
			PathFinder:     fakePathFinder.Spy,
			DataDirManager: fakeDataDirManager,
			Config:         etcdConfig,
		}
	})

	It("can be queried for the port it listens on", func() {
		Expect(etcd.URL()).To(Equal("http://this.is.etcd.listening.for.clients:1234"))
	})

	Context("when given a path to a binary that runs for a long time", func() {
		It("can start and stop that binary", func() {
			sessionBuffer := gbytes.NewBuffer()
			fmt.Fprintf(sessionBuffer, "Everything is dandy")
			fakeSession.BufferReturns(sessionBuffer)

			fakeSession.ExitCodeReturnsOnCall(0, -1)
			fakeSession.ExitCodeReturnsOnCall(1, 143)
			fakePathFinder.ReturnsOnCall(0, "/path/to/some/etcd")

			etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				Expect(command.Path).To(Equal("/path/to/some/etcd"))
				fmt.Fprint(err, "serving insecure client requests on this.is.etcd.listening.for.clients:1234")
				return fakeSession, nil
			}

			By("Starting the Etcd Server")
			err := etcd.Start()
			Expect(err).NotTo(HaveOccurred())

			By("...in turn calling the PathFinder")
			Expect(fakePathFinder.CallCount()).To(Equal(1))
			Expect(fakePathFinder.ArgsForCall(0)).To(Equal("etcd"))

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

			err := etcd.Start()
			Expect(err).To(MatchError(ContainSubstring("Error on directory creation.")))
			Expect(processStarterCounter).To(Equal(0))
		})
	})

	Context("when  the starter returns an error", func() {
		It("propagates the error", func() {
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
