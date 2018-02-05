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

package kinflate

import (
	"encoding/base64"
	"reflect"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/gvkn"
)

func makeUnstructuredEnvConfigMap(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":              name,
				"creationTimestamp": nil,
			},
			"data": map[string]interface{}{
				"DB_USERNAME": "admin",
				"DB_PASSWORD": "somepw",
			},
		},
	}
}

func makeUnstructuredEnvSecret(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":              name,
				"creationTimestamp": nil,
			},
			"type": string(corev1.SecretTypeOpaque),
			"data": map[string]interface{}{
				"DB_USERNAME": base64.StdEncoding.EncodeToString([]byte("admin")),
				"DB_PASSWORD": base64.StdEncoding.EncodeToString([]byte("somepw")),
			},
		},
	}
}

func makeUnstructuredTLSSecret(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":              name,
				"creationTimestamp": nil,
			},
			"type": string(corev1.SecretTypeTLS),
			"data": map[string]interface{}{
				"tls.key": base64.StdEncoding.EncodeToString([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBANLJhPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wo
k/4xIA+ui35/MmNartNuC+BdZ1tMuVCPFZcCAwEAAQJAEJ2N+zsR0Xn8/Q6twa4G
6OB1M1WO+k+ztnX/1SvNeWu8D6GImtupLTYgjZcHufykj09jiHmjHx8u8ZZB/o1N
MQIhAPW+eyZo7ay3lMz1V01WVjNKK9QSn1MJlb06h/LuYv9FAiEA25WPedKgVyCW
SmUwbPw8fnTcpqDWE3yTO3vKcebqMSsCIBF3UmVue8YU3jybC3NxuXq3wNm34R8T
xVLHwDXh/6NJAiEAl2oHGGLz64BuAfjKrqwz7qMYr9HCLIe/YsoWq/olzScCIQDi
D2lWusoe2/nEqfDVVWGWlyJ7yOmqaVm/iNUN9B2N2g==
-----END RSA PRIVATE KEY-----
`)),
				"tls.crt": base64.StdEncoding.EncodeToString([]byte(`-----BEGIN CERTIFICATE-----
MIIB0zCCAX2gAwIBAgIJAI/M7BYjwB+uMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTIwOTEyMjE1MjAyWhcNMTUwOTEyMjE1MjAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJ
hPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wok/4xIA+ui35/MmNa
rtNuC+BdZ1tMuVCPFZcCAwEAAaNQME4wHQYDVR0OBBYEFJvKs8RfJaXTH08W+SGv
zQyKn0H8MB8GA1UdIwQYMBaAFJvKs8RfJaXTH08W+SGvzQyKn0H8MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADQQBJlffJHybjDGxRMqaRmDhX0+6v02TUKZsW
r5QuVbpQhH6u+0UgcW0jp9QwpxoPTLTWGXEWBBBurxFwiCBhkQ+V
-----END CERTIFICATE-----
`)),
			},
		},
	}
}
func TestPopulateMap(t *testing.T) {
	expectedMap := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{
		{
			GVK: schema.GroupVersionKind{
				Version: "v1",
				Kind:    "ConfigMap",
			},
			Name: "envConfigMap",
		}: makeUnstructuredEnvConfigMap("newNameConfigMap"),
		{
			GVK: schema.GroupVersionKind{
				Version: "v1",
				Kind:    "Secret",
			},
			Name: "envSecret",
		}: makeUnstructuredEnvSecret("newNameSecret"),
		{
			GVK: schema.GroupVersionKind{
				Version: "v1",
				Kind:    "Secret",
			},
			Name: "tlsSecret",
		}: makeUnstructuredTLSSecret("newNameTLSSecret"),
	}

	m := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	err := populateMap(m, makeUnstructuredEnvConfigMap("envConfigMap"), "newNameConfigMap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = populateMap(m, makeUnstructuredEnvSecret("envSecret"), "newNameSecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = populateMap(m, makeUnstructuredTLSSecret("tlsSecret"), "newNameTLSSecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(m, expectedMap) {
		t.Fatalf("%#v\ndoesn't match expected\n%#v\n", m, expectedMap)
	}

	err = populateMap(m, makeUnstructuredEnvSecret("envSecret"), "newNameSecret")
	if err == nil || !strings.Contains(err.Error(), "duplicate name") {
		t.Fatalf("expected error to contain %q, but got: %v", "duplicate name", err)
	}
}

func TestPopulateMapOfConfigMapAndSecret(t *testing.T) {
	m := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	manifest := &manifest.Manifest{
		Configmaps: []manifest.ConfigMap{
			{
				Name: "envConfigMap",
				DataSources: manifest.DataSources{
					EnvSource: "examples/simple/instances/exampleinstance/configmap/app.env",
				},
			},
		},
		GenericSecrets: []manifest.GenericSecret{
			{
				Name: "envSecret",
				DataSources: manifest.DataSources{
					EnvSource: "examples/simple/instances/exampleinstance/configmap/app.env",
				},
			},
		},
	}
	expectedMap := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{
		{
			GVK: schema.GroupVersionKind{
				Version: "v1",
				Kind:    "ConfigMap",
			},
			Name: "envConfigMap",
		}: makeUnstructuredEnvConfigMap("envConfigMap-d2c89bt4kk"),
		{
			GVK: schema.GroupVersionKind{
				Version: "v1",
				Kind:    "Secret",
			},
			Name: "envSecret",
		}: makeUnstructuredEnvSecret("envSecret-684h2mm268"),
	}
	err := populateConfigMapAndSecretMap(manifest, m)
	if err != nil {
		t.Fatalf("unexpected erorr: %v", err)
	}
	if !reflect.DeepEqual(m, expectedMap) {
		t.Fatalf("%#v\ndoesn't match expected\n%#v\n", m, expectedMap)
	}
}
