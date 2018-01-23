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

	. "k8s.io/kubectl/pkg/framework/unstructpath"
)

func TestStringSSelectFrom(t *testing.T) {
	s := String().SelectFrom(
		"my string",
		1,
		"your string",
		[]int{3, 4})

	if !reflect.DeepEqual(s, []string{"my string", "your string"}) {
		t.Fatal("SelectFrom should select all integers")
	}
}

func TestStringSFilter(t *testing.T) {
	s := String().
		Filter(StringLength(NumberEqual(4))).
		SelectFrom(
			"one",
			"two",
			"three",
			"four",
			"five")

	if !reflect.DeepEqual(s, []string{"four", "five"}) {
		t.Fatal("SelectFrom should filter selected strings")
	}
}

func TestStringSPredicate(t *testing.T) {
	if !String().Filter(StringLength(NumberEqual(4))).Match("four") {
		t.Fatal("SelectFromor matching element should match")
	}
	if String().Filter(StringLength(NumberEqual(10))).Match("four") {
		t.Fatal("SelectFromor not matching element should not match")
	}
}

func TestStringSFromValueS(t *testing.T) {
	if !Children().String().Filter(StringLength(NumberEqual(4))).Match([]interface{}{"four", "five"}) {
		t.Fatal("SelectFromor should find element that match")
	}
	if Children().String().Filter(StringLength(NumberEqual(4))).Match([]interface{}{"one", "two", "three"}) {
		t.Fatal("SelectFromor shouldn't find element that match")
	}
}
