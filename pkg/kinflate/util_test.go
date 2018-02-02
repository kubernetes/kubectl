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
				Type:       "env",
				NamePrefix: "envConfigMap",
				Generic: manifest.Generic{
					EnvSource: "examples/simple/instances/exampleinstance/configmap/app.env",
				},
			},
		},
		Secrets: []manifest.Secret{
			{
				Type:       "env",
				NamePrefix: "envSecret",
				Generic: manifest.Generic{
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
