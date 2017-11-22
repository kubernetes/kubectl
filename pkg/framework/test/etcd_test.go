package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Etcd", func() {

	Context("when given a path to a binary that runs for a long time", func() {
		It("can start and stop that binary", func() {
			pathToFakeEtcd, err := gexec.Build("k8s.io/kubectl/pkg/framework/test/assets/fakeetcd")
			Expect(err).NotTo(HaveOccurred())
			etcd := Etcd{Path: pathToFakeEtcd}

			By("Starting the Etcd Server")
			session, err := etcd.Start("our etcd url", "our data directory")
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say("Everything is dandy"))
			Expect(session).NotTo(gexec.Exit())

			By("Stopping the Etcd Server")
			session.Terminate()
			Eventually(session).Should(gexec.Exit(143))
		})

	})

	Context("when no path is given", func() {
		It("fails with a helpful error", func() {
			etcd := Etcd{}
			_, err := etcd.Start("our etcd url", "")
			Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
		})
	})

	Context("when given a path to a non-executable", func() {
		It("fails with a helpful error", func() {
			apiServer := Etcd{
				Path: "./etcd.go",
			}
			_, err := apiServer.Start("our etcd url", "")
			Expect(err).To(MatchError(ContainSubstring("./etcd.go: permission denied")))
		})
	})
})
