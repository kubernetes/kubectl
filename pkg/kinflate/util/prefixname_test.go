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

package util

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var prefixNameOps = PrefixNameOptions{
	prefix:      "someprefix-",
	pathConfigs: DefaultNamePrefixPathConfigs,
}

var namePrefixedCm1 = unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name": "someprefix-cm1",
		},
	},
}

var namePrefixedCm2 = unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name": "someprefix-cm2",
		},
	},
}

var namePrefixedM = map[GroupVersionKindName]*unstructured.Unstructured{
	{
		gvk:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
		name: "cm1",
	}: &namePrefixedCm1,
	{
		gvk:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
		name: "cm2",
	}: &namePrefixedCm2,
}

func TestPrefixNameRun(t *testing.T) {
	m := createMap()
	err := prefixNameOps.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(m, namePrefixedM) {
		t.Fatalf("%s doesn't match expected %s", m, namePrefixedM)
	}
}
