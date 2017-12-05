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
			fakeEtcdProcess      *testfakes.FakeFixtureProcess
			fakeAPIServerProcess *testfakes.FakeFixtureProcess
			fakeListenURLGetter  *testfakes.FakeListenURLGetter
			fixtures             Fixtures
		)
		BeforeEach(func() {
			fakeEtcdProcess = &testfakes.FakeFixtureProcess{}
			fakeAPIServerProcess = &testfakes.FakeFixtureProcess{}
			fakeListenURLGetter = &testfakes.FakeListenURLGetter{}
			fixtures = Fixtures{
				Etcd:      fakeEtcdProcess,
				APIServer: fakeAPIServerProcess,
				URLGetter: fakeListenURLGetter.Spy,
			}
		})

		It("can start them", func() {
			err := fixtures.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeListenURLGetter.CallCount()).To(Equal(3))

			By("starting Etcd")
			Expect(fakeEtcdProcess.StartCallCount()).To(Equal(1),
				"the Etcd process should be started exactly once")

			By("starting APIServer")
			Expect(fakeAPIServerProcess.StartCallCount()).To(Equal(1),
				"the APIServer process should be started exactly once")
		})

		Context("when starting etcd fails", func() {
			It("wraps the error", func() {
				fakeEtcdProcess.StartReturns(fmt.Errorf("some error"))
				err := fixtures.Start()
				Expect(fakeListenURLGetter.CallCount()).To(Equal(3))
				Expect(err).To(MatchError(ContainSubstring("some error")))
			})
		})

		Context("when starting APIServer fails", func() {
			It("wraps the error", func() {
				fakeAPIServerProcess.StartReturns(fmt.Errorf("another error"))
				err := fixtures.Start()
				Expect(fakeListenURLGetter.CallCount()).To(Equal(3))
				Expect(err).To(MatchError(ContainSubstring("another error")))
			})
		})

		It("can can clean up the temporary directory and stop", func() {
			fixtures.Stop()
			Expect(fakeEtcdProcess.StopCallCount()).To(Equal(1))
			Expect(fakeAPIServerProcess.StopCallCount()).To(Equal(1))
		})

	})
})
