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
	It("Successfully manages the control plane lifecycle", func() {
		var err error

		controlPlane := test.NewControlPlane()

		By("Starting all the control plane processes")
		err = controlPlane.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected controlPlane to start successfully")

		var apiServerURL, etcdClientURL *url.URL
		etcdClientURL = controlPlane.APIServer.Etcd.URL
		apiServerURL = controlPlane.APIURL()

		isEtcdListeningForClients := isSomethingListeningOnPort(etcdClientURL.Host)
		isAPIServerListening := isSomethingListeningOnPort(apiServerURL.Host)

		By("Ensuring Etcd is listening")
		Expect(isEtcdListeningForClients()).To(BeTrue(),
			fmt.Sprintf("Expected Etcd to listen for clients on %s,", etcdClientURL.Host))

		By("Ensuring APIServer is listening")
		Expect(isAPIServerListening()).To(BeTrue(),
			fmt.Sprintf("Expected APIServer to listen on %s", apiServerURL.Host))

		By("Stopping all the control plane processes")
		err = controlPlane.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected controlPlane to stop successfully")

		By("Ensuring Etcd is not listening anymore")
		Expect(isEtcdListeningForClients()).To(BeFalse(), "Expected Etcd not to listen for clients anymore")

		By("Ensuring APIServer is not listening anymore")
		Expect(isAPIServerListening()).To(BeFalse(), "Expected APIServer not to listen anymore")
	})

	Measure("It should be fast to bring up and tear down the control plane", func(b Benchmarker) {
		b.Time("lifecycle", func() {
			controlPlane := test.NewControlPlane()

			controlPlane.Start()
			controlPlane.Stop()
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
