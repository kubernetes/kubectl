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

var _ = Describe("Apiserver", func() {
	var (
		fakeSession        *testfakes.FakeSimpleSession
		fakeCertDirManager *testfakes.FakeCertDirManager
		apiServer          *APIServer
		apiServerConfig    *APIServerConfig
	)

	BeforeEach(func() {
		fakeSession = &testfakes.FakeSimpleSession{}
		fakeCertDirManager = &testfakes.FakeCertDirManager{}

		apiServerConfig = &APIServerConfig{
			EtcdURL:      "http://this.is.etcd:2345/",
			APIServerURL: "http://this.is.the.API.server:8080",
		}
		apiServer = &APIServer{
			Path:           "",
			CertDirManager: fakeCertDirManager,
			Config:         apiServerConfig,
		}
	})

	Context("when given a path to a binary that runs for a long time", func() {
		It("can start and stop that binary", func() {
			sessionBuffer := gbytes.NewBuffer()
			fmt.Fprint(sessionBuffer, "Everything is fine")
			fakeSession.BufferReturns(sessionBuffer)

			fakeSession.ExitCodeReturnsOnCall(0, -1)
			fakeSession.ExitCodeReturnsOnCall(1, 143)

			apiServer.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				fmt.Fprint(err, "Serving insecurely on this.is.the.API.server:8080")
				return fakeSession, nil
			}

			By("Starting the API Server")
			err := apiServer.Start()
			Expect(err).NotTo(HaveOccurred())

			Eventually(apiServer).Should(gbytes.Say("Everything is fine"))
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(0))
			Expect(apiServer).NotTo(gexec.Exit())
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(1))
			Expect(fakeCertDirManager.CreateCallCount()).To(Equal(1))

			By("Stopping the API Server")
			apiServer.Stop()
			Expect(apiServer).To(gexec.Exit(143))
			Expect(fakeSession.TerminateCallCount()).To(Equal(1))
			Expect(fakeSession.WaitCallCount()).To(Equal(1))
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(2))
			Expect(fakeCertDirManager.DestroyCallCount()).To(Equal(1))
		})
	})

	Context("when the certificate directory cannot be created", func() {
		It("propagates the error, and does not start the process", func() {
			fakeCertDirManager.CreateReturnsOnCall(0, "", fmt.Errorf("Error on cert directory creation."))

			processStarterCounter := 0
			apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				processStarterCounter += 1
				return fakeSession, nil
			}

			err := apiServer.Start()
			Expect(err).To(MatchError(ContainSubstring("Error on cert directory creation.")))
			Expect(processStarterCounter).To(Equal(0))
		})
	})

	Context("when  the starter returns an error", func() {
		It("propagates the error", func() {
			apiServer.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				return nil, fmt.Errorf("Some error in the apiserver starter.")
			}

			err := apiServer.Start()
			Expect(err).To(MatchError(ContainSubstring("Some error in the apiserver starter.")))
		})
	})

	Context("when we try to stop a server that hasn't been started", func() {
		It("is a noop and does not call exit on the session", func() {
			apiServer.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
				return fakeSession, nil
			}
			apiServer.Stop()
			Expect(fakeSession.ExitCodeCallCount()).To(Equal(0))
		})
	})
})
