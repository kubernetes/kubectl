package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultPortFinder", func() {
	It("returns a free port and an address to bind to", func() {
		port, addr, err := DefaultPortFinder("127.0.0.1")

		Expect(err).NotTo(HaveOccurred())
		Expect(addr).To(Equal("127.0.0.1"))
		Expect(port).To(BeNumerically(">=", 1))
		Expect(port).To(BeNumerically("<=", 65535))
	})

	It("errrors on invalid host", func() {
		_, _, err := DefaultPortFinder("this is not a hostname")

		Expect(err).To(MatchError(ContainSubstring("no such host")))
	})

	Context("when using a DNS name", func() {
		It("returns a free port and an address to bind to", func() {
			port, addr, err := DefaultPortFinder("localhost")

			Expect(err).NotTo(HaveOccurred())
			Expect(addr).To(Equal("127.0.0.1"))
			Expect(port).To(BeNumerically(">=", 1))
			Expect(port).To(BeNumerically("<=", 65535))
		})
	})
})
