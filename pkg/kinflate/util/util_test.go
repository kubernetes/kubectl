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

func createMap() map[GroupVersionKindName]*unstructured.Unstructured {
	cm1 := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "cm1",
			},
		},
	}

	cm2 := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "cm2",
			},
		},
	}
	return map[GroupVersionKindName]*unstructured.Unstructured{
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: "cm1",
		}: &cm1,
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: "cm2",
		}: &cm2,
	}
}

func TestDecode(t *testing.T) {
	expected := createMap()
	m, err := Decode(encoded, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%#v doesn't match expected %#v", m, expected)
	}
}

func TestEncode(t *testing.T) {
	out, err := Encode(createMap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(out, encoded) {
		t.Fatalf("%s doesn't match expected %s", out, encoded)
	}
}

func TestFilterByGVK(t *testing.T) {
	type testCase struct {
		description string
		in          schema.GroupVersionKind
		filter      *schema.GroupVersionKind
		expected    bool
	}
	testCases := []testCase{
		{
			description: "nil filter",
			in:          schema.GroupVersionKind{},
			filter:      nil,
			expected:    true,
		},
		{
			description: "GVK matches",
			in: schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			filter: &schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			expected: true,
		},
		{
			description: "group doesn't matches",
			in: schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			filter: &schema.GroupVersionKind{
				Group:   "group2",
				Version: "version1",
				Kind:    "kind1",
			},
			expected: false,
		},
		{
			description: "version doesn't matches",
			in: schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			filter: &schema.GroupVersionKind{
				Group:   "group1",
				Version: "version2",
				Kind:    "kind1",
			},
			expected: false,
		},
		{
			description: "kind doesn't matches",
			in: schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			filter: &schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind2",
			},
			expected: false,
		},
		{
			description: "no version in filter",
			in: schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			filter: &schema.GroupVersionKind{
				Group:   "group1",
				Version: "",
				Kind:    "kind1",
			},
			expected: true,
		},
		{
			description: "only kind is set in filter",
			in: schema.GroupVersionKind{
				Group:   "group1",
				Version: "version1",
				Kind:    "kind1",
			},
			filter: &schema.GroupVersionKind{
				Group:   "",
				Version: "",
				Kind:    "kind1",
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		filtered := SelectByGVK(tc.in, tc.filter)
		if filtered != tc.expected {
			t.Fatalf("unexpected filter result for test case: %v", tc.description)
		}
	}
}

func compareMap(m1, m2 map[GroupVersionKindName]*unstructured.Unstructured) error {
	if len(m1) != len(m2) {
		keySet1 := []GroupVersionKindName{}
		keySet2 := []GroupVersionKindName{}
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
