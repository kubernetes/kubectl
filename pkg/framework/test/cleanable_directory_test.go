package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCertDir", func() {
	It("returns a valid CertDir struct", func() {
		certDir, err := newDirectory()
		Expect(err).NotTo(HaveOccurred())
		Expect(certDir.Path).To(BeADirectory())
		Expect(certDir.Cleanup()).To(Succeed())
		Expect(certDir.Path).NotTo(BeAnExistingFile())
	})
})
