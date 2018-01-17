package internal_test

import (
	"net/url"

	. "k8s.io/kubectl/pkg/framework/test/internal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apiserver", func() {
	It("generates APIServer args", func() {
		processInput := DefaultedProcessInput{
			Dir: "/some/dir",
			URL: url.URL{
				Scheme: "http",
				Host:   "some.apiserver.host:1234",
			},
		}
		etcdURL := &url.URL{
			Scheme: "http",
			Host:   "some.etcd.server",
		}

		args, err := MakeAPIServerArgs(processInput, etcdURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(args)).To(BeNumerically(">", 0))
		Expect(args).To(ContainElement("--etcd-servers=http://some.etcd.server"))
		Expect(args).To(ContainElement("--cert-dir=/some/dir"))
		Expect(args).To(ContainElement("--insecure-port=1234"))
		Expect(args).To(ContainElement("--insecure-bind-address=some.apiserver.host"))
	})

	Context("when URL is not configured", func() {
		It("returns an error", func() {
			var etcdURL *url.URL
			processInput := DefaultedProcessInput{}

			_, err := MakeAPIServerArgs(processInput, etcdURL)
			Expect(err).To(MatchError("must configure Etcd URL"))
		})
	})
})
