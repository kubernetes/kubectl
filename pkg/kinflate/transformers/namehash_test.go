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

func makeSecret(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

func makeHashTestMap() types.KObject {
	return types.KObject{
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: "cm1",
		}: makeConfigmap("cm1"),
		{
			GVK:  schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			Name: "deploy1",
		}: makeDeployment(),
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "Service"},
			Name: "svc1",
		}: makeService(),
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "Secret"},
			Name: "secret1",
		}: makeSecret("secret1"),
	}
}

func makeExpectedHashTestMap() types.KObject {
	return types.KObject{
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
			Name: "cm1",
		}: makeConfigmap("cm1-m462kdfb68"),
		{
			GVK:  schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			Name: "deploy1",
		}: makeDeployment(),
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "Service"},
			Name: "svc1",
		}: makeService(),
		{
			GVK:  schema.GroupVersionKind{Version: "v1", Kind: "Secret"},
			Name: "secret1",
		}: makeSecret("secret1-7kc45hd5f7"),
	}
}

func TestNameHashTransformer(t *testing.T) {
	objs := makeHashTestMap()

	tran := NewNameHashTransformer()
	tran.Transform(objs)

	expected := makeExpectedHashTestMap()

	if !reflect.DeepEqual(objs, expected) {
		err := compareMap(objs, expected)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}
