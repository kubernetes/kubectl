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

package transformers

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/kinflate/types"
)

func makeBaseDeployment() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"old-label": "old-value",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx",
							},
						},
					},
				},
			},
		},
	}
}

func makeOverlayDeployment() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"another-label": "foo",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeMergedDeployment() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"old-label":     "old-value",
							"another-label": "foo",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeTestKObject(genDeployment func() *unstructured.Unstructured) types.KObject {
	return types.KObject{
		{
			GVK:  schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			Name: "deploy1",
		}: genDeployment(),
	}
}

func TestOverlayRun(t *testing.T) {
	base := makeTestKObject(makeBaseDeployment)
	lt, err := NewOverlayTransformer(makeTestKObject(makeOverlayDeployment))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := makeTestKObject(makeMergedDeployment)
	if !reflect.DeepEqual(base, expected) {
		err = compareMap(base, expected)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}
