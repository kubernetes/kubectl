package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Apiserver", func() {

	Context("when given a path to a binary that runs for a long time", func() {
		It("can start and stop that binary", func() {
			pathToFakeAPIServer, err := gexec.Build("k8s.io/kubectl/pkg/framework/test/assets/fakeapiserver")
			Expect(err).NotTo(HaveOccurred())
			apiServer := &APIServer{Path: pathToFakeAPIServer}

			By("Starting the API Server")
			err = apiServer.Start("the etcd url")
			Expect(err).NotTo(HaveOccurred())

			Eventually(apiServer).Should(gbytes.Say("Everything is fine"))
			Expect(apiServer).NotTo(gexec.Exit())

			By("Stopping the API Server")
			apiServer.Stop()
			Eventually(apiServer).Should(gexec.Exit(143))
		})

	})

	Context("when no path is given", func() {
		It("fails with a helpful error", func() {
			apiServer := &APIServer{}
			err := apiServer.Start("the etcd url")
			Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
		})
	})

	Context("when given a path to a non-executable", func() {
		It("fails with a helpful error", func() {
			apiServer := &APIServer{
				Path: "./apiserver.go",
			}
			err := apiServer.Start("the etcd url")
			Expect(err).To(MatchError(ContainSubstring("./apiserver.go: permission denied")))
		})
	})

	Context("when we try to stop a server that hasn't been started", func() {
		It("does not panic", func() {
			server := &APIServer{}
			server.Stop()
		})
	})
})
