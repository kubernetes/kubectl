package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path"
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

	assetsDir := ""

	if dirFromEnv, ok := os.LookupEnv("KUBE_ASSETS_DIR"); ok {
		assetsDir = dirFromEnv
	} else {
		if _, thisFile, _, ok := runtime.Caller(0); ok {
			assetsDir = path.Clean(path.Join(path.Dir(thisFile), "..", "..", "assets", "bin"))
		}
	}

	Expect(assetsDir).NotTo(BeEmpty(),
		"Could not determine assets directory (Hint: you can set $KUBE_ASSETS_DIR)")

	fixtures, err = test.NewFixtures(filepath.Join(assetsDir, "etcd"), filepath.Join(assetsDir, "kube-apiserver"))
	Expect(err).NotTo(HaveOccurred())
	err = fixtures.Start()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	fixtures.Stop()
	gexec.CleanupBuildArtifacts()
})
