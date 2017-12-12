package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewAPIServer", func() {
	It("can construct a properly configured APIServer", func() {
		config := &APIServerConfig{
			APIServerURL: "some APIServer URL",
		}

		apiServer, err := NewAPIServer(config)

		Expect(err).NotTo(HaveOccurred())
		Expect(apiServer).NotTo(BeNil())
		Expect(apiServer.ProcessStarter).NotTo(BeNil())
		Expect(apiServer.CertDirManager).NotTo(BeNil())
		Expect(apiServer.Etcd).NotTo(BeNil())
		Expect(apiServer.Config).To(Equal(config))
	})
})
