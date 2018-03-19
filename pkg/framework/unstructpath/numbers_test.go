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

	p "k8s.io/kubectl/pkg/framework/predicates"
	. "k8s.io/kubectl/pkg/framework/unstructpath"
)

func TestNumberSSelectFrom(t *testing.T) {
	s := Number().SelectFrom(
		1.,
		"string",
		2.,
		[]float64{3, 4})

	if !reflect.DeepEqual(s, []float64{1, 2}) {
		t.Fatal("SelectFrom should select all integers")
	}
}

func TestNumberSFilter(t *testing.T) {
	s := Number().
		Filter(p.NumberGreaterThan(2), p.NumberEqualOrLessThan(4)).
		SelectFrom(
			1.,
			2.,
			3.,
			4.,
			5.)

	if !reflect.DeepEqual(s, []float64{3, 4}) {
		t.Fatal("SelectFrom should filter selected numberegers")
	}
}

func TestNumberSPredicate(t *testing.T) {
	if !Number().Filter(p.NumberGreaterThan(10)).Match(12.) {
		t.Fatal("SelectFromor matching element should match")
	}
	if Number().Filter(p.NumberGreaterThan(10)).Match(4.) {
		t.Fatal("SelectFromor not matching element should not match")
	}
}

func TestNumberSFromValueS(t *testing.T) {
	if !Children().Number().Filter(p.NumberGreaterThan(10)).Match([]interface{}{1., 2., 5., 12.}) {
		t.Fatal("SelectFromor should find element that match")
	}
	if Children().Number().Filter(p.NumberGreaterThan(10)).Match([]interface{}{1., 2., 5.}) {
		t.Fatal("SelectFromor shouldn't find element that match")
	}
}
