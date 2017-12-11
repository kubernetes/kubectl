package test_test

import (
	. "k8s.io/kubectl/pkg/framework/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("APIServerConfig", func() {
	It("does not error on valid config", func() {
		conf := &APIServerConfig{
			APIServerURL: "http://this.is.some.url:1234",
			EtcdURL:      "http://this.is.another.url/with/a/path/we/dont/care/about",
		}
		err := conf.Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("errors on empty config", func() {
		conf := &APIServerConfig{}
		err := conf.Validate()
		Expect(err).To(MatchError(ContainSubstring("APIServerURL: non zero value required")))
	})

	It("errors on malformed URLs", func() {
		conf := &APIServerConfig{
			APIServerURL: "something not URLish",
			EtcdURL:      "something not URLesc",
		}
		err := conf.Validate()
		Expect(err).To(MatchError(ContainSubstring("APIServerURL: something not URLish does not validate as url")))
	})
})
