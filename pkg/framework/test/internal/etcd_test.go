package internal_test

import (
	"net/url"

	. "k8s.io/kubectl/pkg/framework/test/internal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd", func() {
	It("can create Etcd arguments", func() {
		input := DefaultedProcessInput{
			URL: url.URL{
				Scheme: "http",
				Host:   "some.etcd.service:5432",
			},
			Dir: "/some/data/dir",
		}

		args := MakeEtcdArgs(input)
		Expect(args).To(ContainElement("--advertise-client-urls=http://some.etcd.service:5432"))
		Expect(args).To(ContainElement("--listen-client-urls=http://some.etcd.service:5432"))
		Expect(args).To(ContainElement("--data-dir=/some/data/dir"))
	})
})
