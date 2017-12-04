package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"os"
	"path/filepath"

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
	assetsDir, ok := os.LookupEnv("KUBE_ASSETS_DIR")
	Expect(ok).To(BeTrue(), "Expected $KUBE_ASSETS_DIR to be set")

	defaultPathToEtcd = filepath.Join(assetsDir, "etcd")
	defaultPathToApiserver = filepath.Join(assetsDir, "kube-apiserver")
})

var _ = AfterSuite(func() {
	gexec.TerminateAndWait()
})
