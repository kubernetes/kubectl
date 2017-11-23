package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework/test/testfakes"
)

var _ = Describe("Fixtures", func() {
	It("can construct a properly wired Fixtures struct", func() {
		f := NewFixtures("path to etcd", "path to apiserver")
		Expect(f.Etcd.(*Etcd).Path).To(Equal("path to etcd"))
		Expect(f.APIServer.(*APIServer).Path).To(Equal("path to apiserver"))
	})

	Context("with a properly configured set of Fixtures", func() {
		var (
			fakeEtcdStartStopper      *testfakes.FakeEtcdStartStopper
			fakeAPIServerStartStopper *testfakes.FakeAPIServerStartStopper
			fakeTempDirManager        *testfakes.FakeTempDirManager
			fixtures                  Fixtures
		)
		BeforeEach(func() {
			fakeEtcdStartStopper = &testfakes.FakeEtcdStartStopper{}
			fakeAPIServerStartStopper = &testfakes.FakeAPIServerStartStopper{}
			fakeTempDirManager = &testfakes.FakeTempDirManager{}
			fixtures = Fixtures{
				Etcd:           fakeEtcdStartStopper,
				APIServer:      fakeAPIServerStartStopper,
				TempDirManager: fakeTempDirManager,
			}
		})

		It("can start them", func() {
			fakeTempDirManager.CreateReturns("some temp dir")

			err := fixtures.Start()
			Expect(err).NotTo(HaveOccurred())

			By("creating a temporary directory")
			Expect(fakeTempDirManager.CreateCallCount()).To(Equal(1),
				"the TempDirManager should be called exactly once")

			By("starting Etcd")
			Expect(fakeEtcdStartStopper.StartCallCount()).To(Equal(1),
				"the EtcdStartStopper should be called exactly once")
			url, datadir := fakeEtcdStartStopper.StartArgsForCall(0)
			Expect(url).To(Equal("http://127.0.0.1:2379"))
			Expect(datadir).To(Equal("some temp dir"))

			By("starting APIServer")
			Expect(fakeAPIServerStartStopper.StartCallCount()).To(Equal(1),
				"the APIServerStartStopper should be called exactly once")
			url = fakeAPIServerStartStopper.StartArgsForCall(0)
			Expect(url).To(Equal("http://127.0.0.1:2379"))
		})

		Context("when starting etcd fails", func() {
			It("wraps the error", func() {
				fakeEtcdStartStopper.StartReturns(fmt.Errorf("some error"))
				err := fixtures.Start()
				Expect(err).To(MatchError(ContainSubstring("some error")))
			})
		})

		Context("when starting APIServer fails", func() {
			It("wraps the error", func() {
				fakeAPIServerStartStopper.StartReturns(fmt.Errorf("another error"))
				err := fixtures.Start()
				Expect(err).To(MatchError(ContainSubstring("another error")))
			})
		})

		It("can can clean up the temporary directory and stop", func() {
			fixtures.Stop()
			Expect(fakeEtcdStartStopper.StopCallCount()).To(Equal(1))
			Expect(fakeAPIServerStartStopper.StopCallCount()).To(Equal(1))
			Expect(fakeTempDirManager.DestroyCallCount()).To(Equal(1))
		})

		Context("when cleanup fails", func() {
			It("still stops the services, and it bubbles up the error", func() {
				fakeTempDirManager.DestroyReturns(fmt.Errorf("deletion failed"))
				err := fixtures.Stop()
				Expect(err).To(MatchError(ContainSubstring("deletion failed")))

				Expect(fakeEtcdStartStopper.StopCallCount()).To(Equal(1))
				Expect(fakeAPIServerStartStopper.StopCallCount()).To(Equal(1))
			})
		})
	})
})
