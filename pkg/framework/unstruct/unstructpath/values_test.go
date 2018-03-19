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

package unstructpath_test

import (
	"reflect"
	"testing"

	"k8s.io/kubectl/pkg/framework/unstruct"
	"k8s.io/kubectl/pkg/framework/unstruct/unstructpath"
)

func TestRoot(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	roots := unstructpath.Root().Select(
		root,
		root.Map().Field("key1"),
		root.Map().Field("key3").Slice().At(1),
		root.Map().Field("key4").Map().Field("subkey").Slice().At(0),
	)

	if !reflect.DeepEqual(roots, []unstruct.Value{root, root, root, root}) {
		t.Fatal("Root is supposed to alawys return root: ", roots)
	}
}

func TestAll(t *testing.T) {
	u := unstruct.New(map[string]interface{}{
		"key1": 1.,
		"key2": []interface{}{2., 3., map[string]interface{}{"key3": 4.}},
		"key4": map[string]interface{}{"key5": 5.},
	})

	numbers := unstructpath.All().Number().Select(u)
	expected := []float64{1., 2., 3., 4., 5.}
	if !reflect.DeepEqual(expected, numbers) {
		t.Fatalf("Expected to find all numbers (%v), got: %v", expected, numbers)
	}
}

func TestChildren(t *testing.T) {
	u := unstruct.New(map[string]interface{}{
		"key1": 1.,
		"key2": []interface{}{2., 3., map[string]interface{}{"key3": 4.}},
		"key4": 5.,
	})

	numbers := unstructpath.Children().Number().Select(u)
	expected := []float64{1., 5.}
	if !reflect.DeepEqual(expected, numbers) {
		t.Fatalf("Expected to find all numbers (%v), got: %v", expected, numbers)
	}
}

func TestFilter(t *testing.T) {
	us := []unstruct.Value{
		unstruct.New([]interface{}{1., 2., 3.}),
		unstruct.New([]interface{}{3., 4., 5., 6.}),
		unstruct.New(map[string]interface{}{}),
		unstruct.New(5.),
		unstruct.New("string"),
	}
	expected := []unstruct.Value{us[1]}
	actual := unstructpath.Filter(unstructpath.Slice().At(3)).Select(us...)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected to filter (%v), got: %v", expected, actual)
	}
}

func TestParent(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	parents := unstructpath.Parent().Select(
		root, // Root doesn't have parent, shouldn't select anything.
		root.Map().Field("key1"),
		root.Map().Field("key3").Slice().At(1),
		root.Map().Field("key4").Map().Field("subkey").Slice().At(0),
	)
	expected := []unstruct.Value{root, root.Map().Field("key3"), root.Map().Field("key4").Map().Field("subkey")}
	if !reflect.DeepEqual(expected, parents) {
		t.Fatalf("Parent() selected %v, should select %v", parents, expected)
	}
}

func TestValueSPredicate(t *testing.T) {
	if !unstructpath.Slice().Match(unstruct.New([]interface{}{})) {
		t.Fatal("Selecting a slice from a slice should match.")
	}
}

func TestValueSMap(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	expected := []unstruct.Map{
		root.Map(),
		root.Map().Field("key4").Map(),
	}
	actual := unstructpath.All().Map().Select(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Map should find maps %v, got %v", expected, actual)
	}
}

func TestValueSSlice(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	expected := []unstruct.Slice{
		root.Map().Field("key3").Slice(),
		root.Map().Field("key4").Map().Field("subkey").Slice(),
	}
	actual := unstructpath.All().Slice().Select(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Slice should find slices %v, got %v", expected, actual)
	}
}

func TestValueSChildren(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	expected := []unstruct.Value{
		root.Map().Field("key3").Slice().At(0),
		root.Map().Field("key3").Slice().At(1),
	}
	actual := unstructpath.Map().Field("key3").Children().Select(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestValueSNumber(t *testing.T) {
	u := unstruct.New([]interface{}{1., 2., "three", 4., 5., []interface{}{}})

	numbers := unstructpath.Children().Number().Select(u)
	expected := []float64{1., 2., 4., 5.}
	if !reflect.DeepEqual(expected, numbers) {
		t.Fatalf("Children().Number() should select %v, got %v", expected, numbers)

	}
}

func TestValueSString(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	expected := []string{
		"value",
		"other value",
		"string",
	}
	actual := unstructpath.All().String().Select(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestValueSParent(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	expected := []unstruct.Value{
		root,
		root,
		root,
		root,
	}
	// Children().Parent() should return "self" for each children.
	actual := unstructpath.Children().Parent().Select(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestValueSAll(t *testing.T) {
	root := unstruct.New(map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": []interface{}{
			"other value",
			2,
		},
		"key4": map[string]interface{}{
			"subkey": []interface{}{
				3,
				"string",
			},
		},
	})

	expected := []unstruct.Value{
		root.Map().Field("key4"),
		root.Map().Field("key4").Map().Field("subkey"),
		root.Map().Field("key4").Map().Field("subkey").Slice().At(0),
		root.Map().Field("key4").Map().Field("subkey").Slice().At(1),
	}
	// Children().Parent() should return "self" for each children.
	actual := unstructpath.Map().Field("key4").All().Select(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
