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
	"k8s.io/apimachinery/pkg/runtime/schema"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/types"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

func makeMapOfConfigMap() types.KObject {
	return types.KObject{
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

func makeMapOfPod() types.KObject {
	return makeMapOfPodWithImageName("nginx")
}

func makeMapOfPodWithImageName(imageName string) types.KObject {
	return types.KObject{
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
		Resources:  ResourcesType(types.KObject{}),
		Patches:    PatchesType(types.KObject{}),
		Configmaps: ConfigmapsType(types.KObject{}),
		Secrets:    SecretsType(types.KObject{}),
		Packages:   []*ManifestData{},
	}
}

func TestFileToMap(t *testing.T) {
	type testcase struct {
		filename  string
		expected  types.KObject
		expectErr bool
		errorStr  string
	}

	testcases := []testcase{
		{
			filename: "testdata/valid/cm/configmap.yaml",
			expected: types.KObject{
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

	// TODO: Convert to a fake filesystem instead of using test files.
	loader := Loader{FS: fs.MakeRealFS()}

	for _, tc := range testcases {
		actual := types.KObject{}
		err := loader.loadKObjectFromFile(tc.filename, actual)
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
		expected  types.KObject
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

	// TODO: Convert to a fake filesystem instead of using test files.
	loader := Loader{FS: fs.MakeRealFS()}

	for _, tc := range testcases {
		actual := types.KObject{}
		err := loader.loadKObjectFromPath(tc.filename, actual)
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
		expected  types.KObject
		expectErr bool
		errorStr  string
	}

	mapOfConfigMap := makeMapOfConfigMap()
	mapOfPod := makeMapOfPod()
	err := types.Merge(mapOfPod, mapOfConfigMap)
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

	// TODO: Convert to a fake filesystem instead of using test files.
	loader := Loader{FS: fs.MakeRealFS()}

	for _, tc := range testcases {
		actual, err := loader.loadKObjectFromPaths(tc.filenames)
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
	err := types.Merge(mapOfPod, mapOfConfigMap)
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
		Resources:         ResourcesType(mergedMap),
		Patches:           PatchesType(makeMapOfPodWithImageName("nginx:latest")),
		Configmaps:        ConfigmapsType(types.KObject{}),
		Secrets:           SecretsType(types.KObject{}),
	}

	// TODO: Convert to a fake filesystem instead of using test files.
	loader := Loader{FS: fs.MakeRealFS()}
	actual, err := loader.loadManifestDataFromManifestFileAndResources(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expectedMd) {
		t.Errorf("expect:\n%#v\nbut got:\n%#v", expectedMd, actual)
	}
}

func TestLoadManifestDataFromPath(t *testing.T) {
	grandparent := makeManifestData("grandparent")
	parent1 := makeManifestData("parent1")
	parent2 := makeManifestData("parent2")
	child1 := makeManifestData("child1")
	child2 := makeManifestData("child2")
	grandparent.Packages = []*ManifestData{parent1, parent2}
	parent1.Packages = []*ManifestData{child1}
	parent2.Packages = []*ManifestData{child2}

	loader := Loader{FS: fs.MakeRealFS(), InitialPath: "testdata/hierarchy"}
	expected := grandparent
	actual, err := loader.LoadManifestDataFromPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expect:\n%#v\nbut got:\n%#v", expected, actual)
	}
}
