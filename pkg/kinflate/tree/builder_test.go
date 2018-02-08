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

package tree

import (
	"reflect"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/gvkn"
	"k8s.io/kubectl/pkg/kinflate/mergemap"
)

func makeMapOfConfigMap() map[gvkn.GroupVersionKindName]*unstructured.Unstructured {
	return map[gvkn.GroupVersionKindName]*unstructured.Unstructured{
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: "cm1",
		}: {
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm1",
				},
				"data": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}
}

func makeMapOfPod() map[gvkn.GroupVersionKindName]*unstructured.Unstructured {
	return makeMapOfPodWithImageName("nginx")
}

func makeMapOfPodWithImageName(imageName string) map[gvkn.GroupVersionKindName]*unstructured.Unstructured {
	return map[gvkn.GroupVersionKindName]*unstructured.Unstructured{
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "Pod"},
			Name: "pod1",
		}: {
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name": "pod1",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": imageName,
						},
					},
				},
			},
		},
	}
}

func makeManifestData(name string) *ManifestData {
	return &ManifestData{
		Name:       name,
		Resources:  map[gvkn.GroupVersionKindName]*unstructured.Unstructured{},
		Patches:    map[gvkn.GroupVersionKindName]*unstructured.Unstructured{},
		Configmaps: map[gvkn.GroupVersionKindName]*unstructured.Unstructured{},
		Secrets:    map[gvkn.GroupVersionKindName]*unstructured.Unstructured{},
	}
}

func TestValidateManifestPath(t *testing.T) {
	type testcase struct {
		filename  string
		expectErr bool
		errorStr  string
	}

	testcases := []testcase{
		{
			filename:  "testdata/valid/",
			expectErr: false,
		},
		{
			filename:  "testdata/valid/Kube-manifest.yaml",
			expectErr: false,
		},
		{
			filename:  "does-not-exist",
			expectErr: true,
			errorStr:  "no such file or directory",
		},
		{
			filename:  "testdata/invalid/",
			expectErr: true,
			errorStr:  "no such file or directory",
		},
	}

	for _, tc := range testcases {
		_, err := validateManifestPath(tc.filename)
		if err == nil {
			if tc.expectErr {
				t.Errorf("filename: %q, expect an error containing %q, but didn't get an error", tc.filename, tc.errorStr)
			}
		} else {
			if tc.expectErr {
				if !strings.Contains(err.Error(), tc.errorStr) {
					t.Errorf("filename: %q, expect an error containing %q, but got %v", tc.filename, tc.errorStr, err)
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestFileToMap(t *testing.T) {
	type testcase struct {
		filename  string
		expected  map[gvkn.GroupVersionKindName]*unstructured.Unstructured
		expectErr bool
		errorStr  string
	}

	testcases := []testcase{
		{
			filename: "testdata/valid/cm/configmap.yaml",
			expected: map[gvkn.GroupVersionKindName]*unstructured.Unstructured{
				{
					GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
					Name: "cm1",
				}: {
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"foo": "bar",
						},
					},
				},
			},
			expectErr: false,
		},
		{
			filename:  "testdata/valid/cm/",
			expectErr: true,
			errorStr:  "NOT expected to be an dir",
		},
		{
			filename:  "does-not-exist",
			expectErr: true,
			errorStr:  "no such file or directory",
		},
	}

	for _, tc := range testcases {
		actual := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
		err := fileToMap(tc.filename, actual)
		if err == nil {
			if tc.expectErr {
				t.Errorf("filename: %q, expect an error containing %q, but didn't get an error", tc.filename, tc.errorStr)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("filename: %q, expect %v, but got %v", tc.filename, tc.expected, actual)
			}
		} else {
			if tc.expectErr {
				if !strings.Contains(err.Error(), tc.errorStr) {
					t.Errorf("filename: %q, expect an error containing %q, but got %v", tc.filename, tc.errorStr, err)
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestPathToMap(t *testing.T) {
	type testcase struct {
		filename  string
		expected  map[gvkn.GroupVersionKindName]*unstructured.Unstructured
		expectErr bool
		errorStr  string
	}

	expectedMap := makeMapOfConfigMap()

	testcases := []testcase{
		{
			filename:  "testdata/valid/cm/configmap.yaml",
			expected:  expectedMap,
			expectErr: false,
		},
		{
			filename:  "testdata/valid/cm/",
			expected:  expectedMap,
			expectErr: false,
		},
		{
			filename:  "does-not-exist",
			expectErr: true,
			errorStr:  "no such file or directory",
		},
	}

	for _, tc := range testcases {
		actual := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
		err := pathToMap(tc.filename, actual)
		if err == nil {
			if tc.expectErr {
				t.Errorf("filename: %q, expect an error containing %q, but didn't get an error", tc.filename, tc.errorStr)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("filename: %q, expect %v, but got %v", tc.filename, tc.expected, actual)
			}
		} else {
			if tc.expectErr {
				if !strings.Contains(err.Error(), tc.errorStr) {
					t.Errorf("filename: %q, expect an error containing %q, but got %v", tc.filename, tc.errorStr, err)
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestPathsToMap(t *testing.T) {
	type testcase struct {
		filenames []string
		expected  map[gvkn.GroupVersionKindName]*unstructured.Unstructured
		expectErr bool
		errorStr  string
	}

	mapOfConfigMap := makeMapOfConfigMap()
	mapOfPod := makeMapOfPod()
	err := mergemap.Merge(mapOfPod, mapOfConfigMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mergedMap := mapOfPod

	testcases := []testcase{
		{
			filenames: []string{"testdata/valid/cm/"},
			expected:  mapOfConfigMap,
			expectErr: false,
		},
		{
			filenames: []string{"testdata/valid/pod.yaml"},
			expected:  makeMapOfPod(),
			expectErr: false,
		},
		{
			filenames: []string{"testdata/valid/cm/", "testdata/valid/pod.yaml"},
			expected:  mergedMap,
			expectErr: false,
		},
		{
			filenames: []string{"does-not-exist"},
			expectErr: true,
			errorStr:  "no such file or directory",
		},
	}

	for _, tc := range testcases {
		actual, err := pathsToMap(tc.filenames)
		if err == nil {
			if tc.expectErr {
				t.Errorf("filenames: %q, expect an error containing %q, but didn't get an error", tc.filenames, tc.errorStr)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("filenames: %q, expect %v, but got %v", tc.filenames, tc.expected, actual)
			}
		} else {
			if tc.expectErr {
				if !strings.Contains(err.Error(), tc.errorStr) {
					t.Errorf("filenames: %q, expect an error containing %q, but got %v", tc.filenames, tc.errorStr, err)
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestManifestToManifestData(t *testing.T) {
	mapOfConfigMap := makeMapOfConfigMap()
	mapOfPod := makeMapOfPod()
	err := mergemap.Merge(mapOfPod, mapOfConfigMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mergedMap := mapOfPod

	m := &manifest.Manifest{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-manifest",
		},
		NamePrefix: "someprefix-",
		ObjectLabels: map[string]string{
			"foo": "bar",
		},
		ObjectAnnotations: map[string]string{
			"note": "This is an annotation.",
		},
		Resources: []string{
			"testdata/valid/cm/",
			"testdata/valid/pod.yaml",
		},
		Patches: []string{
			"testdata/valid/patch.yaml",
		},
	}

	expectedMd := &ManifestData{
		Name:              "test-manifest",
		NamePrefix:        "someprefix-",
		ObjectLabels:      map[string]string{"foo": "bar"},
		ObjectAnnotations: map[string]string{"note": "This is an annotation."},
		Resources:         mergedMap,
		Patches:           makeMapOfPodWithImageName("nginx:latest"),
		Configmaps:        map[gvkn.GroupVersionKindName]*unstructured.Unstructured{},
		Secrets:           map[gvkn.GroupVersionKindName]*unstructured.Unstructured{},
	}

	actual, err := manifestToManifestData(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expectedMd) {
		t.Errorf("expect:\n%#v\nbut got:\n%#v", expectedMd, actual)
	}
}

func TestManifestPathToManifestNode(t *testing.T) {
	expected := &ManifestNode{
		Data: makeManifestData("grandparent"),
		Children: []*ManifestNode{
			{
				Data: makeManifestData("parent1"),
				Children: []*ManifestNode{
					{
						Data:     makeManifestData("child1"),
						Children: []*ManifestNode{},
					},
				},
			},
			{
				Data: makeManifestData("parent2"),
				Children: []*ManifestNode{
					{
						Data:     makeManifestData("child2"),
						Children: []*ManifestNode{},
					},
				},
			},
		},
	}

	actual, err := manifestPathToManifestNode("testdata/hierarchy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expect:\n%#v\nbut got:\n%#v", expected, actual)
	}
}
