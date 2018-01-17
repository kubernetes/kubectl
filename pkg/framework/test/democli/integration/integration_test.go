package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("DemoCLI Integration", func() {
	It("can give us a helpful help message", func() {
		helpfulMessage := `This is a demo kubernetes CLI, which interacts with the kubernetes API.`

		command := exec.Command(pathToDemoCommand, "--help")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out).To(gbytes.Say(helpfulMessage))
	})

	It("can get a list of pods", func() {
		apiURL := controlPlane.APIURL()

		command := exec.Command(pathToDemoCommand, "listPods", "--api-url", apiURL.String())
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out).To(gbytes.Say("There are no pods."))
	})
})
