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

package unstructpath

import "k8s.io/kubectl/pkg/framework/unstruct"

// ValueS is a "value selector". It filters values based on the
// "filtered" predicates.
type ValueS interface {
	// ValueS can be used as a Value predicate. If the selector can't
	// select any value from the value, then the predicate is
	// false.
	ValueP

	// Select finds values from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	Select(...unstruct.Value) []unstruct.Value

	// Map returns a selector that selects Maps from the given
	// values.
	Map() MapS
	// Slice returns a selector that selects Slices from the given
	// values.
	Slice() SliceS
	// Number returns a selector taht selects Numbers from the given values.
	Number() NumberS
	// String returns a selector that selects strings from the given values.
	String() StringS

	// Parent returns a selector that selects the parent of each
	// given values. If the value is a root, then no value is
	// selected.
	Parent() ValueS
	// Children returns a selector that selects the direct children
	// of the given values.
	Children() ValueS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() ValueS

	// Filter will create a new StringS that filters only the values
	// who match the predicate.
	Filter(...ValueP) ValueS
}

// Root is a ValueS that selects the root of the values.
func Root() ValueS {
	return &valueS{vf: valueRootFilter{}}
}

// Children selects all the children of the values.
func Children() ValueS {
	return &valueS{vf: valueChildrenFilter{}}
}

// All selects all the direct and indirect childrens of the values.
func All() ValueS {
	return &valueS{vf: valueAllFilter{}}
}

// Parent selects the parent of each value. Nothing is selected if the
// value is the root.
func Parent() ValueS {
	return &valueS{vf: valueParentFilter{}}
}

// Filter will only return the values that match the predicate.
func Filter(predicates ...ValueP) ValueS {
	return &valueS{vf: &valueFilterP{vp: ValueAnd(predicates...)}}
}

// ValueS is a "Value Selector". It selects a list of values, maps,
// slices, strings, integer from a list of values.
type valueS struct {
	vs ValueS
	vf valueFilter
}

// Match returns true if the selector can find items in the given
// value. Otherwise, it returns false.
func (s *valueS) Match(value unstruct.Value) bool {
	return len(s.Select(value)) != 0
}

func (s *valueS) Select(values ...unstruct.Value) []unstruct.Value {
	if s.vs != nil {
		values = s.vs.Select(values...)
	}
	return s.vf.Select(values...)
}

func (s *valueS) Map() MapS {
	return &mapS{vs: s}
}

func (s *valueS) Slice() SliceS {
	return &sliceS{vs: s}
}

func (s *valueS) Number() NumberS {
	return &numberS{vs: s}
}

func (s *valueS) String() StringS {
	return &stringS{vs: s}
}

func (s *valueS) Parent() ValueS {
	return &valueS{vs: s, vf: valueParentFilter{}}
}

func (s *valueS) Children() ValueS {
	return &valueS{vs: s, vf: valueChildrenFilter{}}
}

func (s *valueS) All() ValueS {
	return &valueS{vs: s, vf: valueAllFilter{}}
}

func (s *valueS) Filter(predicates ...ValueP) ValueS {
	return &valueS{vs: s, vf: &valueFilterP{vp: ValueAnd(predicates...)}}
}
