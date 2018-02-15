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

	"k8s.io/kubectl/pkg/framework/unstruct"
	. "k8s.io/kubectl/pkg/framework/unstruct/unstructpath"
)

type SliceTrue struct{}

func (SliceTrue) Match(slice unstruct.Slice) bool {
	return true
}

func TestSliceNot(t *testing.T) {
	if SliceNot(SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceNot(SliceTrue{}) should never match")
	}
	if !SliceNot(SliceNot(SliceTrue{})).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceNot(SliceNot(SliceTrue{})) should always match")
	}
}

func TestSliceAnd(t *testing.T) {
	if !SliceAnd().Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceAnd() should always match")
	}
	if SliceAnd(SliceNot(SliceTrue{})).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceAnd(SliceNot(SliceTrue{})) should never match")
	}
	if !SliceAnd(SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceAnd(SliceTrue{}) should always match")
	}
	if !SliceAnd(SliceTrue{}, SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceAnd(SliceTrue{}, SliceTrue{}) should always match")
	}
	if SliceAnd(SliceTrue{}, SliceNot(SliceTrue{}), SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceAnd(SliceTrue{}, SliceNot(SliceTrue{}), SliceTrue{}) should never match")
	}
}

func TestSliceOr(t *testing.T) {
	if SliceOr().Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceOr() should never match")
	}
	if SliceOr(SliceNot(SliceTrue{})).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceOr(SliceNot(SliceTrue{})) should never match")
	}
	if !SliceOr(SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceOr(SliceTrue{}) should always match")
	}
	if !SliceOr(SliceTrue{}, SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceOr(SliceTrue{}, SliceTrue{}) should always match")
	}
	if !SliceOr(SliceTrue{}, SliceNot(SliceTrue{}), SliceTrue{}).Match(unstruct.New(nil).Slice()) {
		t.Fatal("SliceOr(SliceTrue{}, SliceNot(SliceTrue{}), SliceTrue{}) should always match")
	}
}

func TestSliceLength(t *testing.T) {
	slice := unstruct.New([]interface{}{1, 2, 3}).Slice()
	if !SliceLength(NumberEqual(3)).Match(slice) {
		t.Fatal(`SliceLength(NumberEqual(3)) should match []interface{}{1, 2, 3}`)
	}

	if SliceLength(NumberLessThan(2)).Match(slice) {
		t.Fatal(`SliceLength(NumberLessThan(2)) should not match []interface{}{1, 2, 3}`)
	}

	if !SliceLength(NumberLessThan(5)).Match(slice) {
		t.Fatal(`SliceLength(NumberLessThan(5)) should match []interface{}{1, 2, 3}`)
	}
}
