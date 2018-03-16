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
	"testing"

	. "k8s.io/kubectl/pkg/framework/unstructpath"
)

type ValueTrue struct{}

func (ValueTrue) Match(value interface{}) bool {
	return true
}

func TestValueNot(t *testing.T) {
	if ValueNot(ValueTrue{}).Match(nil) {
		t.Fatal("ValueNot(ValueTrue{}) should never match")
	}
	if !ValueNot(ValueNot(ValueTrue{})).Match(nil) {
		t.Fatal("ValueNot(ValueNot(ValueTrue{})) should always match")
	}
}

func TestValueAnd(t *testing.T) {
	if !ValueAnd().Match(nil) {
		t.Fatal("ValueAnd() should always match")
	}
	if ValueAnd(ValueNot(ValueTrue{})).Match(nil) {
		t.Fatal("ValueAnd(ValueNot(ValueTrue{})) should never match")
	}
	if !ValueAnd(ValueTrue{}).Match(nil) {
		t.Fatal("ValueAnd(ValueTrue{}) should always match")
	}
	if !ValueAnd(ValueTrue{}, ValueTrue{}).Match(nil) {
		t.Fatal("ValueAnd(ValueTrue{}, ValueTrue{}) should always match")
	}
	if ValueAnd(ValueTrue{}, ValueNot(ValueTrue{}), ValueTrue{}).Match(nil) {
		t.Fatal("ValueAnd(ValueTrue{}, ValueNot(ValueTrue{}), ValueTrue{}) should never match")
	}
}

func TestValueOr(t *testing.T) {
	if ValueOr().Match(nil) {
		t.Fatal("ValueOr() should never match")
	}
	if ValueOr(ValueNot(ValueTrue{})).Match(nil) {
		t.Fatal("ValueOr(ValueNot(ValueTrue{})) should never match")
	}
	if !ValueOr(ValueTrue{}).Match(nil) {
		t.Fatal("ValueOr(ValueTrue{}) should always match")
	}
	if !ValueOr(ValueTrue{}, ValueTrue{}).Match(nil) {
		t.Fatal("ValueOr(ValueTrue{}, ValueTrue{}) should always match")
	}
	if !ValueOr(ValueTrue{}, ValueNot(ValueTrue{}), ValueTrue{}).Match(nil) {
		t.Fatal("ValueOr(ValueTrue{}, ValueNot(ValueTrue{}), ValueTrue{}) should always match")
	}
}

func TestValueDeepEqual(t *testing.T) {
	if !ValueDeepEqual([]int{1, 2, 3}).Match([]int{1, 2, 3}) {
		t.Fatal("ValueDeepEqual([]int{1, 2, 3}) should match []int{1, 2, 3}")
	}
	if ValueDeepEqual([]int{1, 2, 3}).Match([]int{1, 2, 4}) {
		t.Fatal("ValueDeepEqual([]int{1, 2, 3}) should not match []int{1, 2, 4}")
	}
}
