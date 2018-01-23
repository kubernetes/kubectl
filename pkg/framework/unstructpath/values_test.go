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

	"k8s.io/kubectl/pkg/framework/unstructpath"
)

func TestAll(t *testing.T) {
	u := map[string]interface{}{
		"key1": 1.,
		"key2": []interface{}{2., 3., map[string]interface{}{"key3": 4.}},
		"key4": map[string]interface{}{"key5": 5.},
	}

	numbers := unstructpath.All().Number().SelectFrom(u)
	expected := []float64{1., 2., 3., 4., 5.}
	if !reflect.DeepEqual(expected, numbers) {
		t.Fatalf("Expected to find all numbers (%v), got: %v", expected, numbers)
	}
}

func TestChildren(t *testing.T) {
	u := map[string]interface{}{
		"key1": 1.,
		"key2": []interface{}{2., 3., map[string]interface{}{"key3": 4.}},
		"key4": 5.,
	}

	numbers := unstructpath.Children().Number().SelectFrom(u)
	expected := []float64{1., 5.}
	if !reflect.DeepEqual(expected, numbers) {
		t.Fatalf("Expected to find all numbers (%v), got: %v", expected, numbers)
	}
}

func TestFilter(t *testing.T) {
	us := []interface{}{
		[]interface{}{1., 2., 3.},
		[]interface{}{3., 4., 5., 6.},
		map[string]interface{}{},
		5.,
		"string",
	}
	expected := []interface{}{us[1]}
	actual := unstructpath.Filter(unstructpath.Slice().At(3)).SelectFrom(us...)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected to filter (%v), got: %v", expected, actual)
	}
}

func TestValueSPredicate(t *testing.T) {
	if !unstructpath.Slice().Match([]interface{}{}) {
		t.Fatal("SelectFroming a slice from a slice should match.")
	}
}

func TestValueSMap(t *testing.T) {
	root := map[string]interface{}{
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
	}

	expected := []map[string]interface{}{
		root,
		root["key4"].(map[string]interface{}),
	}
	actual := unstructpath.All().Map().SelectFrom(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Map should find maps %v, got %v", expected, actual)
	}
}

func TestValueSSlice(t *testing.T) {
	root := map[string]interface{}{
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
	}

	expected := [][]interface{}{
		root["key3"].([]interface{}),
		root["key4"].(map[string]interface{})["subkey"].([]interface{}),
	}
	actual := unstructpath.All().Slice().SelectFrom(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Slice should find slices %#v, got %#v", expected, actual)
	}
}

func TestValueSChildren(t *testing.T) {
	root := map[string]interface{}{
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
	}

	expected := []interface{}{
		root["key3"].([]interface{})[0],
		root["key3"].([]interface{})[1],
	}
	actual := unstructpath.Map().Field("key3").Children().SelectFrom(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestValueSNumber(t *testing.T) {
	u := []interface{}{1., 2., "three", 4., 5., []interface{}{}}

	numbers := unstructpath.Children().Number().SelectFrom(u)
	expected := []float64{1., 2., 4., 5.}
	if !reflect.DeepEqual(expected, numbers) {
		t.Fatalf("Children().Number() should select %v, got %v", expected, numbers)

	}
}

func TestValueSString(t *testing.T) {
	root := map[string]interface{}{
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
	}

	expected := []string{
		"value",
		"other value",
		"string",
	}
	actual := unstructpath.All().String().SelectFrom(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestValueSAll(t *testing.T) {
	root := map[string]interface{}{
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
	}

	expected := []interface{}{
		root["key4"],
		root["key4"].(map[string]interface{})["subkey"],
		root["key4"].(map[string]interface{})["subkey"].([]interface{})[0],
		root["key4"].(map[string]interface{})["subkey"].([]interface{})[1],
	}

	actual := unstructpath.Map().Field("key4").All().SelectFrom(root)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
