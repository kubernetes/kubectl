package test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EtcdConfig", func() {
	Context("Constructor", func() {
		var (
			previousPortFinder PortFinder
		)
		BeforeEach(func() {
			previousPortFinder = etcdPortFinder
		})
		AfterEach(func() {
			etcdPortFinder = previousPortFinder
		})

		It("sets some sane default URLs", func() {
			etcdPortFinder = func(host string) (port int, addr string, err error) {
				return 42, "1.2.3.4", nil
			}

			conf, err := NewEtcdConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(conf.ClientURL).To(Equal("http://1.2.3.4:42"))
		})

		It("prop[agates and error while trying to find a port", func() {
			etcdPortFinder = func(host string) (port int, addr string, err error) {
				return 43, "5.6.7.8", fmt.Errorf("oh nooos, no port")
			}

			_, err := NewEtcdConfig()
			Expect(err).To(MatchError(ContainSubstring("oh nooos, no port")))
		})
	})
	Context("Validate()", func() {
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
})
