package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework/test/testfakes"
)

var _ = Describe("ControlPlane", func() {
	Context("with a properly configured set of ControlPlane", func() {
		var (
			fakeAPIServerProcess *testfakes.FakeControlPlaneProcess
			controlPlane         ControlPlane
		)
		BeforeEach(func() {
			fakeAPIServerProcess = &testfakes.FakeControlPlaneProcess{}
			controlPlane = ControlPlane{
				APIServer: fakeAPIServerProcess,
			}
		})

		It("can start them", func() {
			err := controlPlane.Start()
			Expect(err).NotTo(HaveOccurred())

			By("starting APIServer")
			Expect(fakeAPIServerProcess.StartCallCount()).To(Equal(1),
				"the APIServer process should be started exactly once")
		})

		Context("when starting APIServer fails", func() {
			It("propagates the error", func() {
				fakeAPIServerProcess.StartReturns(fmt.Errorf("another error"))
				err := controlPlane.Start()
				Expect(err).To(MatchError(ContainSubstring("another error")))
			})
		})

		It("can can clean up the temporary directory and stop", func() {
			controlPlane.Stop()
			Expect(fakeAPIServerProcess.StopCallCount()).To(Equal(1))
		})

		Context("when stopping APIServer fails", func() {
			It("propagates the error", func() {
				fakeAPIServerProcess.StopReturns(fmt.Errorf("error on stop"))
				err := controlPlane.Stop()
				Expect(err).To(MatchError(ContainSubstring("error on stop")))
			})
		})

		It("can be queried for the APIServer URL", func() {
			fakeAPIServerProcess.URLReturns("some url to the apiserver", nil)

			url, err := controlPlane.APIServerURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(url).To(Equal("some url to the apiserver"))
		})

		Context("when querying the URL fails", func() {
			It("propagates the error", func() {
				fakeAPIServerProcess.URLReturns("", fmt.Errorf("URL error"))
				_, err := controlPlane.APIServerURL()
				Expect(err).To(MatchError(ContainSubstring("URL error")))
			})
		})
	})
})
