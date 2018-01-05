package test_test

import (
	"fmt"
	"io"
	"os/exec"

	"time"

	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "k8s.io/kubectl/pkg/framework/test"
	"k8s.io/kubectl/pkg/framework/test/testfakes"
)

var _ = Describe("Etcd", func() {
	var (
		fakeSession         *testfakes.FakeSimpleSession
		etcd                *Etcd
		etcdStopper         chan struct{}
		dataDirCleanupCount int
	)

	BeforeEach(func() {
		fakeSession = &testfakes.FakeSimpleSession{}

		etcdStopper = make(chan struct{}, 1)
		fakeSession.TerminateReturns(&gexec.Session{
			Exited: etcdStopper,
		})
		close(etcdStopper)

		etcd = &Etcd{
			Path: "/path/to/some/etcd",
			DataDir: &Directory{
				Path: "/path/to/some/etcd",
				Cleanup: func() error {
					dataDirCleanupCount += 1
					return nil
				},
			},
			StopTimeout: 500 * time.Millisecond,
		}
	})

	Describe("starting and stopping etcd", func() {
		Context("when given a path to a binary that runs for a long time", func() {
			It("can start and stop that binary", func() {
				sessionBuffer := gbytes.NewBuffer()
				fmt.Fprintf(sessionBuffer, "Everything is dandy")
				fakeSession.BufferReturns(sessionBuffer)

				fakeSession.ExitCodeReturnsOnCall(0, -1)
				fakeSession.ExitCodeReturnsOnCall(1, 143)

				etcd.Address = &url.URL{Scheme: "http", Host: "this.is.etcd.listening.for.clients:1234"}

				etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					Expect(command.Args).To(ContainElement(fmt.Sprintf("--advertise-client-urls=http://%s:%d", "this.is.etcd.listening.for.clients", 1234)))
					Expect(command.Args).To(ContainElement(fmt.Sprintf("--listen-client-urls=http://%s:%d", "this.is.etcd.listening.for.clients", 1234)))
					Expect(command.Path).To(Equal("/path/to/some/etcd"))
					fmt.Fprint(err, "serving insecure client requests on this.is.etcd.listening.for.clients:1234")
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
				Expect(etcd.Stop()).To(Succeed())

				Expect(dataDirCleanupCount).To(Equal(1))
				Expect(etcd).To(gexec.Exit(143))
				Expect(fakeSession.TerminateCallCount()).To(Equal(1))
				Expect(fakeSession.ExitCodeCallCount()).To(Equal(2))
			})
		})

		Context("when the data directory cannot be destroyed", func() {
			It("propagates the error", func() {
				etcd.DataDir.Cleanup = func() error {
					return fmt.Errorf("destroy failed")
				}
				etcd.Address = &url.URL{Scheme: "http", Host: "this.is.etcd:1234"}
				etcd.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					fmt.Fprint(err, "serving insecure client requests on this.is.etcd:1234")
					return fakeSession, nil
				}

				Expect(etcd.Start()).To(Succeed())
				err := etcd.Stop()
				Expect(err).To(MatchError(ContainSubstring("destroy failed")))
			})
		})

		Context("when there is no function to cleanup the data directory", func() {
			It("does not panic", func() {
				etcd.DataDir.Cleanup = nil
				etcd.Address = &url.URL{Scheme: "http", Host: "this.is.etcd:1234"}
				etcd.ProcessStarter = func(Command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					fmt.Fprint(err, "serving insecure client requests on this.is.etcd:1234")
					return fakeSession, nil
				}

				Expect(etcd.Start()).To(Succeed())

				var err error
				Expect(func() {
					err = etcd.Stop()
				}).NotTo(Panic())
				Expect(err).NotTo(HaveOccurred())
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

		Context("when the starter takes longer than our timeout", func() {
			It("gives us a timeout error", func() {
				etcd.StartTimeout = 1 * time.Nanosecond
				etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					return &gexec.Session{}, nil
				}

				err := etcd.Start()
				Expect(err).To(MatchError(ContainSubstring("timeout waiting for etcd to start serving")))
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

		Context("when Stop() times out", func() {
			JustBeforeEach(func() {
				etcdStopperWillNotBeUsed := make(chan struct{})
				fakeSession.TerminateReturns(&gexec.Session{
					Exited: etcdStopperWillNotBeUsed,
				})
			})
			It("propagates the error", func() {
				etcd.Address = &url.URL{Scheme: "http", Host: "this.is.etcd:1234"}

				etcd.ProcessStarter = func(command *exec.Cmd, out, err io.Writer) (SimpleSession, error) {
					fmt.Fprint(err, "serving insecure client requests on this.is.etcd:1234")
					return fakeSession, nil
				}

				Expect(etcd.Start()).To(Succeed())
				err := etcd.Stop()
				Expect(err).To(MatchError(ContainSubstring("timeout")))
			})
		})
	})

	Describe("querying the server for its URL", func() {
		It("can be queried for the URL it listens on", func() {
			etcd.Address = &url.URL{Scheme: "http", Host: "the.host.for.etcd:6789"}
			apiServerURL, err := etcd.URL()
			Expect(err).NotTo(HaveOccurred())
			Expect(apiServerURL).To(Equal("http://the.host.for.etcd:6789"))
		})
		Context("before starting etcd", func() {
			Context("and therefore the addressmanager has not been initialized", func() {
				BeforeEach(func() {
					etcd = &Etcd{}
				})
				It("gives a sane error", func() {
					_, err := etcd.URL()
					Expect(err).To(MatchError(ContainSubstring("not initialized")))
				})
			})
		})
	})
})
