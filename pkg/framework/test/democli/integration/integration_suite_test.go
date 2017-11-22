package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"os"

	"path/filepath"

	"github.com/onsi/gomega/gexec"
	"k8s.io/kubectl/pkg/framework/test"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var (
	pathToDemoCommand string
	fixtures          *test.Fixtures
)

var _ = BeforeSuite(func() {
	var err error
	pathToDemoCommand, err = gexec.Build("k8s.io/kubectl/pkg/framework/test/democli/")
	Expect(err).NotTo(HaveOccurred())

	assetsDir, ok := os.LookupEnv("KUBE_ASSETS_DIR")
	Expect(ok).To(BeTrue(), "KUBE_ASSETS_DIR should point to a directory containing etcd and apiserver binaries")
	fixtures = test.NewFixtures(filepath.Join(assetsDir, "etcd"), filepath.Join(assetsDir, "kube-apiserver"))
	err = fixtures.Start()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	fixtures.Stop()
	gexec.CleanupBuildArtifacts()
})
