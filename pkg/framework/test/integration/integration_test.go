package integration_test

import (
	"fmt"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework/test"
)

var _ = Describe("The Testing Framework", func() {
	It("Successfully manages the fixtures lifecycle", func() {
		fixtures := test.NewFixtures(defaultPathToEtcd, defaultPathToApiserver)

		err := fixtures.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to start successfully")

		isEtcdListening := isSomethingListeningOnPort(2379)
		isAPIServerListening := isSomethingListeningOnPort(8080)

		Expect(isEtcdListening()).To(BeTrue(), "Expected Etcd to listen on 2379")

		Expect(isAPIServerListening()).To(BeTrue(), "Expected APIServer to listen on 8080")

		err = fixtures.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to stop successfully")

		Expect(isEtcdListening()).To(BeFalse(), "Expected Etcd not to listen anymore")

		By("Ensuring APIServer is not listening anymore")
		Expect(isAPIServerListening()).To(BeFalse(), "Expected APIServer not to listen anymore")
	})

	Measure("It should be fast to bring up and tear down the fixtures", func(b Benchmarker) {
		b.Time("lifecycle", func() {
			fixtures := test.NewFixtures(defaultPathToEtcd, defaultPathToApiserver)

			fixtures.Start()
			fixtures.Stop()
		})
	}, 10)
})

type portChecker func() bool

func isSomethingListeningOnPort(port int) portChecker {
	return func() bool {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 1*time.Second)

		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
}
