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
		var err error
		var fixtures *test.Fixtures

		fixtures, err = test.NewFixtures(defaultPathToEtcd, defaultPathToApiserver)
		Expect(err).NotTo(HaveOccurred())

		By("Starting all the fixture processes")
		err = fixtures.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to start successfully")

		apiServerConf := fixtures.APIServer.(*test.APIServer).Config
		etcdConf := fixtures.Etcd.(*test.Etcd).Config

		var apiServerURL, etcdClientURL, etcdPeerURL *url.URL
		etcdClientURL, err = url.Parse(etcdConf.ClientURL)
		Expect(err).NotTo(HaveOccurred())
		etcdPeerURL, err = url.Parse(etcdConf.PeerURL)
		Expect(err).NotTo(HaveOccurred())
		apiServerURL, err = url.Parse(apiServerConf.APIServerURL)
		Expect(err).NotTo(HaveOccurred())

		isEtcdListeningForClients := isSomethingListeningOnPort(etcdClientURL.Host)
		isEtcdListeningForPeers := isSomethingListeningOnPort(etcdPeerURL.Host)
		isAPIServerListening := isSomethingListeningOnPort(apiServerURL.Host)

		By("Ensuring Etcd is listening")
		Expect(isEtcdListeningForClients()).To(BeTrue(),
			fmt.Sprintf("Expected Etcd to listen for clients on %s,", etcdClientURL.Host))
		Expect(isEtcdListeningForPeers()).To(BeTrue(),
			fmt.Sprintf("Expected Etcd to listen for peers on %s,", etcdPeerURL.Host))

		By("Ensuring APIServer is listening")
		Expect(isAPIServerListening()).To(BeTrue(),
			fmt.Sprintf("Expected APIServer to listen on %s", apiServerURL.Host))

		By("Stopping all the fixture processes")
		err = fixtures.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected fixtures to stop successfully")

		By("Ensuring Etcd is not listening anymore")
		Expect(isEtcdListeningForClients()).To(BeFalse(), "Expected Etcd not to listen for clients anymore")
		Expect(isEtcdListeningForPeers()).To(BeFalse(), "Expected Etcd not to listen for peers anymore")

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
