package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"
	"runtime"
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

	_, thisFile, _, ok := runtime.Caller(0)
	Expect(ok).NotTo(BeFalse())
	defaultAssetsDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "assets", "bin"))
	pathToEtcd := filepath.Join(defaultAssetsDir, "etcd")
	pathToAPIServer := filepath.Join(defaultAssetsDir, "kube-apiserver")

	if pathToBin, ok := os.LookupEnv("TEST_ETCD_BIN"); ok {
		pathToEtcd = pathToBin
	}
	if pathToBin, ok := os.LookupEnv("TEST_APISERVER_BIN"); ok {
		pathToAPIServer = pathToBin
	}

	Expect(pathToEtcd).NotTo(BeEmpty(), "Path to etcd cannot be empty, set $TEST_ETCD_BIN")
	Expect(pathToAPIServer).NotTo(BeEmpty(), "Path to apiserver cannot be empty, set $TEST_APISERVER_BIN")

	fixtures, err = test.NewFixtures(pathToEtcd, pathToAPIServer)
	Expect(err).NotTo(HaveOccurred())

	err = fixtures.Start()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	fixtures.Stop()
	gexec.CleanupBuildArtifacts()
})
