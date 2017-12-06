package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EtcdConfig", func() {
	It("does not error on valid config", func() {
		conf := &EtcdConfig{
			PeerURL:   "http://this.is.some.url:1234",
			ClientURL: "http://this.is.another.url/",
		}
		err := conf.Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("errors on empty config", func() {
		conf := &EtcdConfig{}
		err := conf.Validate()
		Expect(err).To(MatchError(ContainSubstring("PeerURL: non zero value required")))
		Expect(err).To(MatchError(ContainSubstring("ClientURL: non zero value required")))
	})

	It("errors on malformed URLs", func() {
		conf := &EtcdConfig{
			PeerURL:   "something not URLish",
			ClientURL: "something not URLesc",
		}
		err := conf.Validate()
		Expect(err).To(MatchError(ContainSubstring("PeerURL: something not URLish does not validate as url")))
		Expect(err).To(MatchError(ContainSubstring("ClientURL: something not URLesc does not validate as url")))
	})
})
