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

import (
	p "k8s.io/kubectl/pkg/framework/predicates"
)

// ValueS is a "value selector". It filters values based on the
// "filtered" predicates.
type ValueS interface {
	// ValueS can be used as a Value predicate. If the selector can't
	// select any value from the value, then the predicate is
	// false.
	p.Value

	// SelectFrom finds values from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) []interface{}

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

	// Children returns a selector that selects the direct children
	// of the given values.
	Children() ValueS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() ValueS

	// Filter will create a new StringS that filters only the values
	// who match the predicate.
	Filter(...p.Value) ValueS
}

// Children selects all the children of the values.
func Children() ValueS {
	return &valueS{vf: valueChildrenFilter{}}
}

// All selects all the direct and indirect childrens of the values.
func All() ValueS {
	return &valueS{vf: valueAllFilter{}}
}

// Filter will only return the values that match the predicate.
func Filter(predicates ...p.Value) ValueS {
	return &valueS{vf: &valueFilterP{vp: p.ValueAnd(predicates...)}}
}

// ValueS is a "Value SelectFromor". It selects a list of values, maps,
// slices, strings, integer from a list of values.
type valueS struct {
	vs ValueS
	vf valueFilter
}

// Match returns true if the selector can find items in the given
// value. Otherwise, it returns false.
func (s *valueS) Match(value interface{}) bool {
	return len(s.SelectFrom(value)) != 0
}

func (s *valueS) SelectFrom(values ...interface{}) []interface{} {
	if s.vs != nil {
		values = s.vs.SelectFrom(values...)
	}
	return s.vf.SelectFrom(values...)
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

func (s *valueS) Children() ValueS {
	return &valueS{vs: s, vf: valueChildrenFilter{}}
}

func (s *valueS) All() ValueS {
	return &valueS{vs: s, vf: valueAllFilter{}}
}

func (s *valueS) Filter(predicates ...p.Value) ValueS {
	return &valueS{vs: s, vf: &valueFilterP{vp: p.ValueAnd(predicates...)}}
}
