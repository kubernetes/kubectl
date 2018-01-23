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
	. "k8s.io/kubectl/pkg/framework/unstruct/unstructpath"
)

func TestNumberSSelect(t *testing.T) {
	s := Number().Select(
		unstruct.New(1.),
		unstruct.New("string"),
		unstruct.New(2.),
		unstruct.New([]float64{3, 4}))

	if !reflect.DeepEqual(s, []float64{1, 2}) {
		t.Fatal("Select should select all integers")
	}
}

func TestNumberSFilter(t *testing.T) {
	s := Number().
		Filter(NumberGreaterThan(2), NumberEqualOrLessThan(4)).
		Select(
			unstruct.New(1.),
			unstruct.New(2.),
			unstruct.New(3.),
			unstruct.New(4.),
			unstruct.New(5.))

	if !reflect.DeepEqual(s, []float64{3, 4}) {
		t.Fatal("Select should filter selected numberegers")
	}
}

func TestNumberSPredicate(t *testing.T) {
	if !Number().Filter(NumberGreaterThan(10)).Match(unstruct.New(12.)) {
		t.Fatal("Selector matching element should match")
	}
	if Number().Filter(NumberGreaterThan(10)).Match(unstruct.New(4.)) {
		t.Fatal("Selector not matching element should not match")
	}
}

func TestNumberSFromValueS(t *testing.T) {
	if !Children().Number().Filter(NumberGreaterThan(10)).Match(unstruct.New([]interface{}{1., 2., 5., 12.})) {
		t.Fatal("Selector should find element that match")
	}
	if Children().Number().Filter(NumberGreaterThan(10)).Match(unstruct.New([]interface{}{1., 2., 5.})) {
		t.Fatal("Selector shouldn't find element that match")
	}
}
