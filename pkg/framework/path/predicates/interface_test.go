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

package predicates_test

import (
	"testing"

	. "k8s.io/kubectl/pkg/framework/path/predicates"
)

type InterfaceTrue struct{}

func (InterfaceTrue) Match(value interface{}) bool {
	return true
}

func TestInterfaceNot(t *testing.T) {
	if InterfaceNot(InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceNot(InterfaceTrue{}) should never match")
	}
	if !InterfaceNot(InterfaceNot(InterfaceTrue{})).Match(nil) {
		t.Fatal("InterfaceNot(InterfaceNot(InterfaceTrue{})) should always match")
	}
}

func TestInterfaceAnd(t *testing.T) {
	if !InterfaceAnd().Match(nil) {
		t.Fatal("InterfaceAnd() should always match")
	}
	if InterfaceAnd(InterfaceNot(InterfaceTrue{})).Match(nil) {
		t.Fatal("InterfaceAnd(InterfaceNot(InterfaceTrue{})) should never match")
	}
	if !InterfaceAnd(InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceAnd(InterfaceTrue{}) should always match")
	}
	if !InterfaceAnd(InterfaceTrue{}, InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceAnd(InterfaceTrue{}, InterfaceTrue{}) should always match")
	}
	if InterfaceAnd(InterfaceTrue{}, InterfaceNot(InterfaceTrue{}), InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceAnd(InterfaceTrue{}, InterfaceNot(InterfaceTrue{}), InterfaceTrue{}) should never match")
	}
}

func TestInterfaceOr(t *testing.T) {
	if InterfaceOr().Match(nil) {
		t.Fatal("InterfaceOr() should never match")
	}
	if InterfaceOr(InterfaceNot(InterfaceTrue{})).Match(nil) {
		t.Fatal("InterfaceOr(InterfaceNot(InterfaceTrue{})) should never match")
	}
	if !InterfaceOr(InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceOr(InterfaceTrue{}) should always match")
	}
	if !InterfaceOr(InterfaceTrue{}, InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceOr(InterfaceTrue{}, InterfaceTrue{}) should always match")
	}
	if !InterfaceOr(InterfaceTrue{}, InterfaceNot(InterfaceTrue{}), InterfaceTrue{}).Match(nil) {
		t.Fatal("InterfaceOr(InterfaceTrue{}, InterfaceNot(InterfaceTrue{}), InterfaceTrue{}) should always match")
	}
}

func TestInterfaceDeepEqual(t *testing.T) {
	if !InterfaceDeepEqual([]int{1, 2, 3}).Match([]int{1, 2, 3}) {
		t.Fatal("InterfaceDeepEqual([]int{1, 2, 3}) should match []int{1, 2, 3}")
	}
	if InterfaceDeepEqual([]int{1, 2, 3}).Match([]int{1, 2, 4}) {
		t.Fatal("InterfaceDeepEqual([]int{1, 2, 3}) should not match []int{1, 2, 4}")
	}
}
