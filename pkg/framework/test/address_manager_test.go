package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultAddressManager", func() {
	var defaultAddressManager *DefaultAddressManager
	BeforeEach(func() {
		defaultAddressManager = &DefaultAddressManager{}
	})

	Describe("Initialize", func() {
		It("returns a free port and an address to bind to", func() {
			port, host, err := defaultAddressManager.Initialize("localhost")

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

		Context("when given an invalid hostname", func() {
			It("propagates the error", func() {
				_, _, err := defaultAddressManager.Initialize("this is not a hostname")

				Expect(err).To(MatchError(ContainSubstring("no such host")))
			})
		})

		Context("when given a hostname that we don't have permission to listen on", func() {
			It("propagates the error", func() {
				_, _, err := defaultAddressManager.Initialize("example.com")

				Expect(err).To(MatchError(ContainSubstring("bind: can't assign requested address")))
			})
		})

		Context("initialized multiple times", func() {
			It("fails", func() {
				_, _, err := defaultAddressManager.Initialize("localhost")
				Expect(err).NotTo(HaveOccurred())
				_, _, err = defaultAddressManager.Initialize("localhost")
				Expect(err).To(MatchError(ContainSubstring("already initialized")))
			})
		})
	})
	Describe("Port", func() {
		It("returns an error if Initialize has not been called yet", func() {
			_, err := defaultAddressManager.Port()
			Expect(err).To(MatchError(ContainSubstring("not initialized yet")))
		})
		It("returns the same port as previously allocated by Initialize", func() {
			expectedPort, _, err := defaultAddressManager.Initialize("localhost")
			Expect(err).NotTo(HaveOccurred())
			actualPort, err := defaultAddressManager.Port()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualPort).To(Equal(expectedPort))
		})
	})
	Describe("Host", func() {
		It("returns an error if Initialize has not been called yet", func() {
			_, err := defaultAddressManager.Host()
			Expect(err).To(MatchError(ContainSubstring("not initialized yet")))
		})
		It("returns the same port as previously allocated by Initialize", func() {
			_, expectedHost, err := defaultAddressManager.Initialize("localhost")
			Expect(err).NotTo(HaveOccurred())
			actualHost, err := defaultAddressManager.Host()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualHost).To(Equal(expectedHost))
		})

	})
})
