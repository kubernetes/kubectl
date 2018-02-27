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
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/kinflate/types"
)

var encoded = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm2
`)

func makeConfigMap(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

func makeConfigMaps(name1InGVKN, name2InGVKN, name1InObj, name2InObj string) types.KObject {
	cm1 := makeConfigMap(name1InObj)
	cm2 := makeConfigMap(name2InObj)
	return types.KObject{
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: name1InGVKN,
		}: cm1,
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: name2InGVKN,
		}: cm2,
	}
}

func TestDecodeToKObject(t *testing.T) {
	expected := makeConfigMaps("cm1", "cm2", "cm1", "cm2")
	m, err := DecodeToKObject(encoded, nil)
	fmt.Printf("%v\n", m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%#v doesn't match expected %#v", m, expected)
	}
}

func TestEncodeFromKObject(t *testing.T) {
	out, err := EncodeFromKObject(makeConfigMaps("cm1", "cm2", "cm1", "cm2"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(out, encoded) {
		t.Fatalf("%s doesn't match expected %s", out, encoded)
	}
}

func compareMap(m1, m2 types.KObject) error {
	if len(m1) != len(m2) {
		keySet1 := []types.GroupVersionKindName{}
		keySet2 := []types.GroupVersionKindName{}
		for GVKn := range m1 {
			keySet1 = append(keySet1, GVKn)
		}
		for GVKn := range m1 {
			keySet2 = append(keySet2, GVKn)
		}
		return fmt.Errorf("maps has different number of entries: %#v doesn't equals %#v", keySet1, keySet2)
	}
	for GVKn, obj1 := range m1 {
		obj2, found := m2[GVKn]
		if !found {
			return fmt.Errorf("%#v doesn't exist in %#v", GVKn, m2)
		}
		if !reflect.DeepEqual(obj1, obj2) {
			return fmt.Errorf("%#v doesn't match %#v", obj1, obj2)
		}
	}
	return nil
}
