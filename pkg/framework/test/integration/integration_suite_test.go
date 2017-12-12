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
	defaultPathToEtcd      string
	defaultPathToApiserver string
)

var _ = BeforeSuite(func() {
	_, thisFile, _, ok := runtime.Caller(0)
	Expect(ok).NotTo(BeFalse())
	defaultAssetsDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "assets", "bin"))
	defaultPathToEtcd = filepath.Join(defaultAssetsDir, "etcd")
	defaultPathToApiserver = filepath.Join(defaultAssetsDir, "kube-apiserver")

	if pathToBin, ok := os.LookupEnv("TEST_ETCD_BIN"); ok {
		defaultPathToEtcd = pathToBin
	}
	if pathToBin, ok := os.LookupEnv("TEST_APISERVER_BIN"); ok {
		defaultPathToApiserver = pathToBin
	}

	Expect(defaultPathToEtcd).NotTo(BeEmpty(), "Path to etcd cannot be empty, set $TEST_ETCD_BIN")
	Expect(defaultPathToApiserver).NotTo(BeEmpty(), "Path to apiserver cannot be empty, set $TEST_APISERVER_BIN")
})

var _ = AfterSuite(func() {
	gexec.TerminateAndWait()
})
