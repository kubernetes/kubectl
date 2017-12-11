package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewAPIServer", func() {
	var oldAPIServerBinPathFinder BinPathFinder
	BeforeEach(func() {
		oldAPIServerBinPathFinder = apiServerBinPathFinder
	})
	AfterEach(func() {
		apiServerBinPathFinder = oldAPIServerBinPathFinder
	})

	It("can construct a properly configured APIServer", func() {
		config := &APIServerConfig{
			EtcdURL:      "some etcd URL",
			APIServerURL: "some APIServer URL",
		}
		apiServerBinPathFinder = func(name string) string {
			Expect(name).To(Equal("kube-apiserver"))
			return "some api server path"
		}

		apiServer, err := NewAPIServer(config)

		Expect(err).NotTo(HaveOccurred())
		Expect(apiServer).NotTo(BeNil())
		Expect(apiServer.ProcessStarter).NotTo(BeNil())
		Expect(apiServer.CertDirManager).NotTo(BeNil())
		Expect(apiServer.Path).To(Equal("some api server path"))
		Expect(apiServer.Etcd).NotTo(BeNil())
		Expect(apiServer.Config).To(Equal(config))
	})
})
