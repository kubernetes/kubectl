package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
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
	defaultPathToEtcd string
)

var _ = BeforeSuite(func() {
	_, thisFile, _, ok := runtime.Caller(0)
	Expect(ok).NotTo(BeFalse())
	defaultAssetsDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "assets", "bin"))
	defaultPathToEtcd = filepath.Join(defaultAssetsDir, "etcd")

	if pathToBin, ok := os.LookupEnv("TEST_ETCD_BIN"); ok {
		defaultPathToEtcd = pathToBin
	}

	Expect(defaultPathToEtcd).NotTo(BeEmpty(), "Path to etcd cannot be empty, set $TEST_ETCD_BIN")
})

var _ = AfterSuite(func() {
	gexec.TerminateAndWait()
})
