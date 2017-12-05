package integration_test

import (
	"net"
	"time"

	"net/url"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework/test"
)

var _ = Describe("The Testing Framework", func() {
	It("Successfully manages the fixtures lifecycle", func() {
		fixtures, err := test.NewFixtures(defaultPathToEtcd, defaultPathToApiserver)
		Expect(err).NotTo(HaveOccurred())

		By("Starting all the fixture processes")
		err = fixtures.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to start successfully")

		var etcdURL, etcdPeerURL, apiServerURL *url.URL
		etcd := fixtures.Etcd.(*test.Etcd)
		apiServer := fixtures.APIServer.(*test.APIServer)

		etcdURL, err = url.Parse(etcd.EtcdURL)
		Expect(err).NotTo(HaveOccurred())
		etcdPeerURL, err = url.Parse(etcd.EtcdPeerURL)
		Expect(err).NotTo(HaveOccurred())
		apiServerURL, err = url.Parse(apiServer.APIServerURL)
		Expect(err).NotTo(HaveOccurred())

		isEtcdListening := isSomethingListeningOnPort(etcdURL.Host)
		isEtcdPeerListening := isSomethingListeningOnPort(etcdPeerURL.Host)
		isAPIServerListening := isSomethingListeningOnPort(apiServerURL.Host)

		By("Ensuring Etcd is listening")
		Expect(isEtcdListening()).To(BeTrue(),
			fmt.Sprintf("Expected Etcd to listen on %s", etcdURL.Host))
		Expect(isEtcdPeerListening()).To(BeTrue(),
			fmt.Sprintf("Expected Etcd to listen for peers on %s", etcdPeerURL.Host))

		By("Ensuring APIServer is listening")
		Expect(isAPIServerListening()).To(BeTrue(),
			fmt.Sprintf("Expected APIServer to listen on %s", apiServerURL.Host))

		By("Stopping all the fixture processes")
		err = fixtures.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to stop successfully")

		By("Ensuring Etcd is not listening anymore")
		Expect(isEtcdListening()).To(BeFalse(), "Expected Etcd not to listen anymore")
		Expect(isEtcdPeerListening()).To(BeFalse(), "Expected Etcd not to listen for peers anymore")

		By("Ensuring APIServer is not listening anymore")
		Expect(isAPIServerListening()).To(BeFalse(), "Expected APIServer not to listen anymore")
	})

	Measure("It should be fast to bring up and tear down the fixtures", func(b Benchmarker) {
		b.Time("lifecycle", func() {
			fixtures, err := test.NewFixtures(defaultPathToEtcd, defaultPathToApiserver)
			Expect(err).NotTo(HaveOccurred())

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
