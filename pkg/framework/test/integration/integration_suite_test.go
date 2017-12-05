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
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Framework Integration Suite")
}

var (
	defaultPathToEtcd      string
	defaultPathToApiserver string
)

var _ = BeforeSuite(func() {
	assetsDir := ""

	if dirFromEnv, ok := os.LookupEnv("KUBE_ASSETS_DIR"); ok {
		assetsDir = dirFromEnv
	} else {
		if _, thisFile, _, ok := runtime.Caller(0); ok {
			assetsDir = path.Clean(path.Join(path.Dir(thisFile), "..", "assets", "bin"))
		}
	}

	Expect(assetsDir).NotTo(BeEmpty(),
		"Could not determine assets directory (Hint: you can set $KUBE_ASSETS_DIR)")

	defaultPathToEtcd = filepath.Join(assetsDir, "etcd")
	defaultPathToApiserver = filepath.Join(assetsDir, "kube-apiserver")
})

var _ = AfterSuite(func() {
	gexec.TerminateAndWait()
})
