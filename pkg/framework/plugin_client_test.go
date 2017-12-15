package framework_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/kubectl/pkg/framework"
)

var _ = Describe("InitConfig", func() {
	var ()

	BeforeEach(func() {
		os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG", "assets/config")
	})

	Describe("InitConfig", func() {
		Context("When nothing is overridden by the calling framework", func() {
			It("finds and parses the preexisting config", func() {
				config, err := framework.InitConfig()
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Host).To(Equal("https://notreal.com:1234"))
				Expect(config.Username).To(Equal("foo"))
				Expect(config.Password).To(Equal("bar"))
			})
		})

		Context("When the calling plugin framework sets env vars", func() {
			BeforeEach(func() {
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_AS", "apple")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_AS_GROUP", "[\"banana\",\"cherry\"]")

				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CERTIFICATE_AUTHORITY", "assets/apiserver_ca.crt")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CLIENT_CERTIFICATE", "assets/client.crt")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CLIENT_KEY", "assets/client.key")

				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_REQUEST_TIMEOUT", "45s")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_SERVER", "some-other-server.com")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_TOKEN", "bearer notreal")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_USERNAME", "date")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_PASSWORD", "elderberry")

				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CLUSTER", "")
			})
			It("overrides the config settings with the passed in settings", func() {
				config, err := framework.InitConfig()
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Impersonate.UserName).To(Equal("apple"))
				Expect(config.Impersonate.Groups).Should(ConsistOf("banana", "cherry"))

				Expect(config.CertFile).To(Equal("assets/client.crt"))
				Expect(config.KeyFile).To(Equal("assets/client.key"))
				Expect(config.CAFile).To(Equal("assets/apiserver_ca.crt"))

				Expect(config.Timeout).To(Equal(45 * time.Second))
				Expect(config.ServerName).To(Equal("some-other-server.com"))
				Expect(config.BearerToken).To(Equal("bearer notreal"))

				Expect(config.Username).To(Equal("date"))
				Expect(config.Password).To(Equal("elderberry"))
			})
		})
	})
})
