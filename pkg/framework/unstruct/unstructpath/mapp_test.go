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

type MapTrue struct{}

func (MapTrue) Match(m unstruct.Map) bool {
	return true
}

func TestMapNot(t *testing.T) {
	if MapNot(MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapNot(MapTrue{}) should never match")
	}
	if !MapNot(MapNot(MapTrue{})).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapNot(MapNot(MapTrue{})) should always match")
	}
}

func TestMapAnd(t *testing.T) {
	if !MapAnd().Match(unstruct.New(nil).Map()) {
		t.Fatal("MapAnd() should always match")
	}
	if MapAnd(MapNot(MapTrue{})).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapAnd(MapNot(MapTrue{})) should never match")
	}
	if !MapAnd(MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapAnd(MapTrue{}) should always match")
	}
	if !MapAnd(MapTrue{}, MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapAnd(MapTrue{}, MapTrue{}) should always match")
	}
	if MapAnd(MapTrue{}, MapNot(MapTrue{}), MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapAnd(MapTrue{}, MapNot(MapTrue{}), MapTrue{}) should never match")
	}
}

func TestMapOr(t *testing.T) {
	if MapOr().Match(unstruct.New(nil).Map()) {
		t.Fatal("MapOr() should never match")
	}
	if MapOr(MapNot(MapTrue{})).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapOr(MapNot(MapTrue{})) should never match")
	}
	if !MapOr(MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapOr(MapTrue{}) should always match")
	}
	if !MapOr(MapTrue{}, MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapOr(MapTrue{}, MapTrue{}) should always match")
	}
	if !MapOr(MapTrue{}, MapNot(MapTrue{}), MapTrue{}).Match(unstruct.New(nil).Map()) {
		t.Fatal("MapOr(MapTrue{}, MapNot(MapTrue{}), MapTrue{}) should always match")
	}
}

func TestMapNumFields(t *testing.T) {
	m := unstruct.New(map[string]interface{}{"First": 1, "Second": 2, "Third": 3}).Map()
	if !MapNumFields(NumberEqual(3)).Match(m) {
		t.Fatal(`MapNumFields(NumberEqual(3)) should match []interface{}{1, 2, 3}`)
	}

	if MapNumFields(NumberLessThan(2)).Match(m) {
		t.Fatal(`MapNumFields(NumberLessThan(2)) should not match []interface{}{1, 2, 3}`)
	}

	if !MapNumFields(NumberLessThan(5)).Match(m) {
		t.Fatal(`MapNumFields(NumberLessThan(5)) should match []interface{}{1, 2, 3}`)
	}
}
