/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"net/http"
	"path/filepath"
	"sync"

	"github.com/golang/protobuf/proto"
	openapiv2 "github.com/googleapis/gnostic/OpenAPIv2"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/framework/openapi"
	"k8s.io/kubectl/pkg/kinflate/transport"
	"k8s.io/kubectl/pkg/pluginutils"
)

const (
	FlagHTTPCacheDir = "cache-dir"
	// protobuf mime type
	mimePb = "application/com.github.proto-openapi.spec.v2@v1.0+protobuf"
)

type factory struct {
	cacheDir     string
	flags        *pflag.FlagSet
	clientConfig *restclient.Config
	// openAPIGetter loads and caches openapi specs
	openAPIGetter openAPIGetter
}

type openAPIGetter struct {
	once   sync.Once
	getter openapi.Getter
}

func newFactory() *factory {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)

	clientConfig, err := pluginutils.InitConfig()
	if err != nil {
		panic(err)
	}

	return &factory{flags: flags, clientConfig: clientConfig}
}

func (f *factory) bindFlags(flags *pflag.FlagSet) {
	// Merge factory's flags
	flags.AddFlagSet(f.flags)

	defaultCacheDir := filepath.Join(homedir.HomeDir(), ".kube", "http-cache")
	// flags.StringVar(&f.cacheDir, FlagHTTPCacheDir, defaultCacheDir, "Default HTTP cache directory")
	flags.StringVar(&f.cacheDir, FlagHTTPCacheDir, defaultCacheDir, "Default HTTP cache directory")
}

// OpenAPISchema returns metadata and structural information about Kubernetes object definitions.
func (f *factory) openAPISchema() (openapi.Resources, error) {
	openAPIClient, err := f.openAPIClient()
	if err != nil {
		return nil, err
	}

	// Lazily initialize the OpenAPIGetter once
	f.openAPIGetter.once.Do(func() {
		// Create the caching OpenAPIGetter
		f.openAPIGetter.getter = openapi.NewOpenAPIGetter(openAPIClient)
	})

	// Delegate to the OpenAPIGetter
	return f.openAPIGetter.getter.Get()
}

func (f *factory) openAPIClient() (*openAPIClient, error) {
	cfg := f.clientConfig

	if f.cacheDir != "" {
		wt := cfg.WrapTransport
		cfg.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			if wt != nil {
				rt = wt(rt)
			}
			return transport.NewCacheRoundTripper(f.cacheDir, rt)
		}
	}

	if err := setDiscoveryDefaults(cfg); err != nil {
		return nil, err
	}
	client, err := restclient.UnversionedRESTClientFor(cfg)
	return &openAPIClient{restClient: client}, err
}

func setDiscoveryDefaults(config *restclient.Config) error {
	config.APIPath = ""
	config.GroupVersion = nil
	codec := runtime.NoopEncoder{Decoder: scheme.Codecs.UniversalDecoder()}
	config.NegotiatedSerializer = serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: codec})
	if len(config.UserAgent) == 0 {
		config.UserAgent = restclient.DefaultKubernetesUserAgent()
	}
	return nil
}

type openAPIClient struct {
	restClient *restclient.RESTClient
	discovery.OpenAPISchemaInterface
}

// OpenAPISchema fetches the open api schema using a rest client and parses the proto.
func (d *openAPIClient) OpenAPISchema() (*openapiv2.Document, error) {
	data, err := d.restClient.Get().AbsPath("/openapi/v2").SetHeader("Accept", mimePb).Do().Raw()
	if err != nil {
		if errors.IsForbidden(err) || errors.IsNotFound(err) {
			// single endpoint not found/registered in old server, try to fetch old endpoint
			// TODO(roycaihw): remove this in 1.11
			data, err = d.restClient.Get().AbsPath("/swagger-2.0.0.pb-v1").Do().Raw()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	document := &openapiv2.Document{}
	err = proto.Unmarshal(data, document)
	if err != nil {
		return nil, err
	}
	return document, nil
}
