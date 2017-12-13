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
		fakeEtcdProcess    *testfakes.FakeFixtureProcess
		fakePathFinder     *testfakes.FakeBinPathFinder
		fakeAddressManager *testfakes.FakeAddressManager
	)

	BeforeEach(func() {
		fakeSession = &testfakes.FakeSimpleSession{}
		fakeCertDirManager = &testfakes.FakeCertDirManager{}
		fakeEtcdProcess = &testfakes.FakeFixtureProcess{}
		fakePathFinder = &testfakes.FakeBinPathFinder{}
		fakeAddressManager = &testfakes.FakeAddressManager{}

		apiServer = &APIServer{
			AddressManager: fakeAddressManager,
			PathFinder:     fakePathFinder.Spy,
			CertDirManager: fakeCertDirManager,
			Etcd:           fakeEtcdProcess,
		}
	})

	Describe("starting and stopping the server", func() {
		Context("when given a path to a binary that runs for a long time", func() {
			It("can start and stop that binary", func() {
				//TODO Break this long test up
				sessionBuffer := gbytes.NewBuffer()
				fmt.Fprint(sessionBuffer, "Everything is fine")
				fakeSession.BufferReturns(sessionBuffer)

				fakeSession.ExitCodeReturnsOnCall(0, -1)
				fakeSession.ExitCodeReturnsOnCall(1, 143)

				fakePathFinder.Returns("/some/path/to/apiserver")
				fakeAddressManager.InitializeReturns(1234, "this.is.the.API.server", nil)
				fakeEtcdProcess.URLReturns("the etcd url", nil)

				apiServer.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					Expect(command.Args).To(ContainElement("--insecure-port=1234"))
					Expect(command.Args).To(ContainElement("--insecure-bind-address=this.is.the.API.server"))
					Expect(command.Args).To(ContainElement("--etcd-servers=the etcd url"))
					Expect(command.Path).To(Equal("/some/path/to/apiserver"))
					fmt.Fprint(err, "Serving insecurely on this.is.the.API.server:1234")
					return fakeSession, nil
				}

				By("Starting the API Server")
				err := apiServer.Start()
				Expect(err).NotTo(HaveOccurred())

				By("...in turn starting Etcd")
				Expect(fakeEtcdProcess.StartCallCount()).To(Equal(1),
					"the Etcd process should be started exactly once")

				By("...in turn calling the PathFinder")
				Expect(fakePathFinder.CallCount()).To(Equal(1))
				Expect(fakePathFinder.ArgsForCall(0)).To(Equal("kube-apiserver"))

				By("...in turn calling the AddressManager")
				Expect(fakeAddressManager.InitializeCallCount()).To(Equal(1))
				Expect(fakeAddressManager.InitializeArgsForCall(0)).To(Equal("localhost"))

				By("...in turn calling the CertDirManager")
				Expect(fakeCertDirManager.CreateCallCount()).To(Equal(1))

				By("...getting the URL of Etcd")
				Expect(fakeEtcdProcess.URLCallCount()).To(Equal(1))

				Eventually(apiServer).Should(gbytes.Say("Everything is fine"))
				Expect(fakeSession.ExitCodeCallCount()).To(Equal(0))
				Expect(apiServer).NotTo(gexec.Exit())
				Expect(fakeSession.ExitCodeCallCount()).To(Equal(1))
				Expect(fakeCertDirManager.CreateCallCount()).To(Equal(1))

				By("Stopping the API Server")
				apiServer.Stop()

				Expect(fakeCertDirManager.DestroyCallCount()).To(Equal(1))
				Expect(fakeEtcdProcess.StopCallCount()).To(Equal(1))
				Expect(apiServer).To(gexec.Exit(143))
				Expect(fakeSession.TerminateCallCount()).To(Equal(1))
				Expect(fakeSession.WaitCallCount()).To(Equal(1))
				Expect(fakeSession.ExitCodeCallCount()).To(Equal(2))
				Expect(fakeCertDirManager.DestroyCallCount()).To(Equal(1))
			})
		})

		Context("when the certificate directory cannot be destroyed", func() {
			It("propagates the error", func() {
				fakeCertDirManager.DestroyReturns(fmt.Errorf("destroy failed"))
				fakeAddressManager.InitializeReturns(1234, "this.is.apiserver", nil)
				apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					fmt.Fprint(err, "Serving insecurely on this.is.apiserver:1234")
					return fakeSession, nil
				}

				Expect(apiServer.Start()).To(Succeed())
				err := apiServer.Stop()
				Expect(err).To(MatchError(ContainSubstring("destroy failed")))
			})
		})

		Context("when etcd cannot be stopped", func() {
			It("propagates the error", func() {
				fakeEtcdProcess.StopReturns(fmt.Errorf("stopping etcd failed"))
				fakeAddressManager.InitializeReturns(1234, "this.is.apiserver", nil)
				apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					fmt.Fprint(err, "Serving insecurely on this.is.apiserver:1234")
					return fakeSession, nil
				}

				Expect(apiServer.Start()).To(Succeed())
				err := apiServer.Stop()
				Expect(err).To(MatchError(ContainSubstring("stopping etcd failed")))
			})
		})

		Context("when starting etcd fails", func() {
			It("propagates the error, and does not start the process", func() {
				fakeEtcdProcess.StartReturnsOnCall(0, fmt.Errorf("starting etcd failed"))
				apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					Expect(true).To(BeFalse(),
						"the api server process starter shouldn't be called if starting etcd fails")
					return nil, nil
				}

				err := apiServer.Start()
				Expect(err).To(MatchError(ContainSubstring("starting etcd failed")))
			})
		})

		Context("when getting the URL of Etcd fails", func() {
			It("propagates the error, stop Etcd and keep APIServer down", func() {
				fakeEtcdProcess.URLReturns("", fmt.Errorf("no etcd url"))

				apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					Expect(true).To(BeFalse(),
						"the api server process starter shouldn't be called if getting etcd's URL fails")
					return nil, nil
				}

				err := apiServer.Start()
				Expect(err).To(MatchError(ContainSubstring("no etcd url")))
				Expect(fakeEtcdProcess.StopCallCount()).To(Equal(1))
			})

			Context("and stopping of etcd fails too", func() {
				It("propagates the combined error", func() {
					fakeEtcdProcess.URLReturns("", fmt.Errorf("no etcd url"))
					fakeEtcdProcess.StopReturns(fmt.Errorf("stopping etcd failed"))

					apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
						Expect(true).To(BeFalse(),
							"the api server process starter shouldn't be called if getting etcd's URL fails")
						return nil, nil
					}

					err := apiServer.Start()
					Expect(err).To(MatchError(ContainSubstring("no etcd url")))
					Expect(err).To(MatchError(ContainSubstring("stopping etcd failed")))
					Expect(fakeEtcdProcess.StopCallCount()).To(Equal(1))
				})
			})
		})

		Context("when the certificate directory cannot be created", func() {
			It("propagates the error, and does not start any process", func() {
				fakeCertDirManager.CreateReturnsOnCall(0, "", fmt.Errorf("Error on cert directory creation."))

				apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					Expect(true).To(BeFalse(),
						"the api server process starter shouldn't be called if creating the cert dir fails")
					return nil, nil
				}

				err := apiServer.Start()
				Expect(err).To(MatchError(ContainSubstring("Error on cert directory creation.")))
				Expect(fakeEtcdProcess.StartCallCount()).To(Equal(0))
			})
		})

		Context("when the address manager fails to get a new address", func() {
			It("propagates the error and does not start any process", func() {
				fakeAddressManager.InitializeReturns(0, "", fmt.Errorf("some error finding a free port"))

				apiServer.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					Expect(true).To(BeFalse(),
						"the api server process starter shouldn't be called if getting a free port fails")
					return nil, nil
				}

				Expect(apiServer.Start()).To(MatchError(ContainSubstring("some error finding a free port")))
				Expect(fakeEtcdProcess.StartCallCount()).To(Equal(0))
			})
		})

		Context("when the starter returns an error", func() {
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

	Describe("querying the server for its URL", func() {
		It("can be queried for the URL it listens on", func() {
			fakeAddressManager.HostReturns("the.host.for.api.server", nil)
			fakeAddressManager.PortReturns(5678, nil)
			apiServerURL, err := apiServer.URL()
			Expect(err).NotTo(HaveOccurred())
			Expect(apiServerURL).To(Equal("http://the.host.for.api.server:5678"))
		})

		Context("when we query for the URL before starting the server", func() {
			Context("and so the addressmanager fails to give us a port", func() {
				It("propagates the failure", func() {
					fakeAddressManager.PortReturns(0, fmt.Errorf("boom"))
					_, err := apiServer.URL()
					Expect(err).To(MatchError(ContainSubstring("boom")))
				})
			})
			Context("and so the addressmanager fails to give us a host", func() {
				It("propagates the failure", func() {
					fakeAddressManager.HostReturns("", fmt.Errorf("biff!"))
					_, err := apiServer.URL()
					Expect(err).To(MatchError(ContainSubstring("biff!")))
				})
			})
		})
	})
})
