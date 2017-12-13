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
		_, err := NewFixtures()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("with a properly configured set of Fixtures", func() {
		var (
			fakeAPIServerProcess *testfakes.FakeFixtureProcess
			fixtures             Fixtures
		)
		BeforeEach(func() {
			fakeAPIServerProcess = &testfakes.FakeFixtureProcess{}
			fixtures = Fixtures{
				APIServer: fakeAPIServerProcess,
			}
		})

		It("can start them", func() {
			err := fixtures.Start()
			Expect(err).NotTo(HaveOccurred())

			By("starting APIServer")
			Expect(fakeAPIServerProcess.StartCallCount()).To(Equal(1),
				"the APIServer process should be started exactly once")
		})

		Context("when starting APIServer fails", func() {
			It("propagates the error", func() {
				fakeAPIServerProcess.StartReturns(fmt.Errorf("another error"))
				err := fixtures.Start()
				Expect(err).To(MatchError(ContainSubstring("another error")))
			})
		})

		It("can can clean up the temporary directory and stop", func() {
			fixtures.Stop()
			Expect(fakeAPIServerProcess.StopCallCount()).To(Equal(1))
		})

		Context("when stopping APIServer fails", func() {
			It("propagates the error", func() {
				fakeAPIServerProcess.StopReturns(fmt.Errorf("error on stop"))
				err := fixtures.Stop()
				Expect(err).To(MatchError(ContainSubstring("error on stop")))
			})
		})

		It("can be queried for the APIServer URL", func() {
			fakeAPIServerProcess.URLReturns("some url to the apiserver", nil)

			url, err := fixtures.APIServerURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(url).To(Equal("some url to the apiserver"))
		})

		Context("when querying the URL fails", func() {
			It("propagates the error", func() {
				fakeAPIServerProcess.URLReturns("", fmt.Errorf("URL error"))
				_, err := fixtures.APIServerURL()
				Expect(err).To(MatchError(ContainSubstring("URL error")))
			})
		})
	})
})
