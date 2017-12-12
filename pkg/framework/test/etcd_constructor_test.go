package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd", func() {
	Context("when constructed with the zero-config constructor", func() {
		var (
			previousBinPathFinder BinPathFinder
		)
		BeforeEach(func() {
			previousBinPathFinder = etcdBinPathFinder
			etcdBinPathFinder = func(name string) (binPath string) {
				return "/the/path/to/etcd"
			}
		})
		AfterEach(func() {
			etcdBinPathFinder = previousBinPathFinder
		})
		It("gets a sensible path", func() {
			etcd, err := NewEtcd()
			Expect(err).NotTo(HaveOccurred())
			Expect(etcd.Path).To(Equal("/the/path/to/etcd"))
		})
	})
})
