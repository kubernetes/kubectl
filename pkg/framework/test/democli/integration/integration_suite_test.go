package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/onsi/gomega/gexec"
	"k8s.io/kubectl/pkg/framework/test"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DemoCLI Integration Suite")
}

var (
	pathToDemoCommand string
	fixtures          *test.Fixtures
)

var _ = BeforeSuite(func() {
	var err error
	pathToDemoCommand, err = gexec.Build("k8s.io/kubectl/pkg/framework/test/democli/")
	Expect(err).NotTo(HaveOccurred())

	fixtures, err = test.NewFixtures()
	Expect(err).NotTo(HaveOccurred())

	err = fixtures.Start()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	fixtures.Stop()
	gexec.CleanupBuildArtifacts()
})
