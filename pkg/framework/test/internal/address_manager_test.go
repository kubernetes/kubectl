package internal_test

import (
	. "k8s.io/kubectl/pkg/framework/test/internal"

	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressManager", func() {
	var addressManager *AddressManager
	BeforeEach(func() {
		addressManager = &AddressManager{}
	})

	Describe("Initialize", func() {
		It("returns a free port and an address to bind to", func() {
			port, host, err := addressManager.Initialize()

			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("127.0.0.1"))
			Expect(port).NotTo(Equal(0))

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
			Expect(err).NotTo(HaveOccurred())
			l, err := net.ListenTCP("tcp", addr)
			defer func() {
				Expect(l.Close()).To(Succeed())
			}()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("initialized multiple times", func() {
			It("fails", func() {
				_, _, err := addressManager.Initialize()
				Expect(err).NotTo(HaveOccurred())
				_, _, err = addressManager.Initialize()
				Expect(err).To(MatchError(ContainSubstring("already initialized")))
			})
		})
	})
	Describe("Port", func() {
		It("returns an error if Initialize has not been called yet", func() {
			_, err := addressManager.Port()
			Expect(err).To(MatchError(ContainSubstring("not initialized yet")))
		})
		It("returns the same port as previously allocated by Initialize", func() {
			expectedPort, _, err := addressManager.Initialize()
			Expect(err).NotTo(HaveOccurred())
			actualPort, err := addressManager.Port()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualPort).To(Equal(expectedPort))
		})
	})
	Describe("Host", func() {
		It("returns an error if Initialize has not been called yet", func() {
			_, err := addressManager.Host()
			Expect(err).To(MatchError(ContainSubstring("not initialized yet")))
		})
		It("returns the same port as previously allocated by Initialize", func() {
			_, expectedHost, err := addressManager.Initialize()
			Expect(err).NotTo(HaveOccurred())
			actualHost, err := addressManager.Host()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualHost).To(Equal(expectedHost))
		})
	})
})
