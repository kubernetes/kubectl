package integration_test

import (
	"net"
	"time"

	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework/test"
)

var _ = Describe("The Testing Framework", func() {
	It("Successfully manages the fixtures lifecycle", func() {
		fixtures := test.NewFixtures(defaultPathToEtcd, defaultPathToApiserver)

		By("Starting all the fixture processes")
		err := fixtures.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to start successfully")

		apiServerURL, err := url.Parse(fixtures.Config.APIServerURL)
		Expect(err).NotTo(HaveOccurred())

		isAPIServerListening := isSomethingListeningOnPort(apiServerURL.Host)

		By("Ensuring APIServer is listening")
		Expect(isAPIServerListening()).To(BeTrue(),
			fmt.Sprintf("Expected APIServer to listen on %s", apiServerURL.Host))

		By("Stopping all the fixture processes")
		err = fixtures.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to stop successfully")

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

func isSomethingListeningOnPort(hostAndPort string) portChecker {
	return func() bool {
		conn, err := net.DialTimeout("tcp", hostAndPort, 1*time.Second)

		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
}
