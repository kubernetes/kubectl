package pluginutils_test

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"time"

	"k8s.io/kubectl/pkg/pluginutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plugin client", func() {
	BeforeEach(func() {
		os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG", "testdata/config")
	})

	Describe("InitClientAndConfig", func() {
		Context("When nothing is overridden by the calling framework", func() {
			It("finds and parses the preexisting config", func() {
				client, config, err := pluginutils.InitClientAndConfig()
				Expect(err).NotTo(HaveOccurred())

				Expect(client.Host).To(Equal("https://notreal.com:1234"))
				Expect(client.Username).To(Equal("foo"))
				Expect(client.Password).To(Equal("bar"))

				namespace, overridden, err := config.Namespace()
				Expect(err).NotTo(HaveOccurred())
				Expect(namespace).To(Equal("default"))
				Expect(overridden).To(BeFalse())
			})
		})

		Context("When the calling plugin framework sets env vars", func() {
			BeforeEach(func() {
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_AS", "apple")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_AS_GROUP", "[\"banana\",\"cherry\"]")

				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CERTIFICATE_AUTHORITY", "testdata/apiserver_ca.crt")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CLIENT_CERTIFICATE", "testdata/client.crt")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CLIENT_KEY", "testdata/client.key")

				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_REQUEST_TIMEOUT", "45s")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_SERVER", "some-other-server.com")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_TOKEN", "bearer notreal")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_USERNAME", "date")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_PASSWORD", "elderberry")

				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_CONTEXT", "california")
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_NAMESPACE", "catalog")
			})
			It("overrides the config settings with the passed in settings", func() {
				client, config, err := pluginutils.InitClientAndConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(client.Impersonate.UserName).To(Equal("apple"))
				Expect(client.Impersonate.Groups).Should(ConsistOf("banana", "cherry"))

				Expect(client.CertFile).To(Equal("testdata/client.crt"))
				Expect(client.KeyFile).To(Equal("testdata/client.key"))
				Expect(client.CAFile).To(Equal("testdata/apiserver_ca.crt"))

				Expect(client.Timeout).To(Equal(45 * time.Second))
				Expect(client.ServerName).To(Equal("some-other-server.com"))
				Expect(client.BearerToken).To(Equal("bearer notreal"))

				Expect(client.Username).To(Equal("date"))
				Expect(client.Password).To(Equal("elderberry"))

				Expect(client.Host).To(Equal("https://notrealincalifornia.com:1234"))

				namespace, overridden, err := config.Namespace()
				Expect(err).NotTo(HaveOccurred())
				Expect(namespace).To(Equal("catalog"))
				Expect(overridden).To(BeTrue())
			})
		})
	})

	Describe("InitClientAndConfig in Base64", func() {
		Context("When nothing is overridden by the calling framework", func() {
			BeforeEach(func() {
				path := os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG")
				b, _ := ioutil.ReadFile(path)
				sEnc := base64.StdEncoding.EncodeToString(b)
				os.Setenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG", "base64:"+sEnc)
			})

			It("just makes sure that the base64 config works", func() {
				client, _, err := pluginutils.InitClientAndConfig()
				Expect(err).NotTo(HaveOccurred())

				Expect(client.Host).To(Equal("https://notrealincalifornia.com:1234"))
				Expect(client.Username).To(Equal("date"))
				Expect(client.Password).To(Equal("elderberry"))
			})
		})
	})
})
