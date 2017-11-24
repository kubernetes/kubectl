package integration_test

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework/test"
)

var _ = Describe("Integration", func() {
	It("Successfully manages the fixtures lifecycle", func() {
		assetsDir, ok := os.LookupEnv("KUBE_ASSETS_DIR")
		Expect(ok).To(BeTrue(), "Expected $KUBE_ASSETS_DIR to be set")

		pathToEtcd := filepath.Join(assetsDir, "etcd")
		pathToApiserver := filepath.Join(assetsDir, "kube-apiserver")

		fixtures := test.NewFixtures(pathToEtcd, pathToApiserver)

		err := fixtures.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to start successfully")

		Eventually(func() bool {
			return isSomethingListeningOnPort(2379)
		}, 5*time.Second).Should(BeTrue(), "Expected Etcd to listen on 2379")

		Eventually(func() bool {
			return isSomethingListeningOnPort(8080)
		}, 5*time.Second).Should(BeTrue(), "Expected APIServer to listen on 8080")

		err = fixtures.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to stop successfully")

		Expect(isSomethingListeningOnPort(2379)).To(BeFalse(), "Expected Etcd not to listen anymore")

		By("Ensuring APIServer is not listening anymore")
		Expect(isSomethingListeningOnPort(8080)).To(BeFalse(), "Expected APIServer not to listen anymore")
	})
})

func isSomethingListeningOnPort(port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
