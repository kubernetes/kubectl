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
			fixtures                  Fixtures
		)
		BeforeEach(func() {
			fakeEtcdStartStopper = &testfakes.FakeEtcdStartStopper{}
			fakeAPIServerStartStopper = &testfakes.FakeAPIServerStartStopper{}
			fixtures = Fixtures{
				Etcd:      fakeEtcdStartStopper,
				APIServer: fakeAPIServerStartStopper,
			}
		})

		It("can start them", func() {
			err := fixtures.Start()
			Expect(err).NotTo(HaveOccurred())

			By("starting Etcd")
			Expect(fakeEtcdStartStopper.StartCallCount()).To(Equal(1),
				"the EtcdStartStopper should be called exactly once")

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
		})

	})
})
