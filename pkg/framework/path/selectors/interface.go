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

package selectors

import (
	p "k8s.io/kubectl/pkg/framework/path/predicates"
)

// Interface is a "interface selector". It filters interfaces based on the
// "filtered" predicates.
type Interface interface {
	// Interface can be used as a Interface predicate. If the selector can't
	// select any interface from the interface, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds interfaces from interfaces using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) []interface{}

	// AsMap returns a selector that selects Maps from the given
	// values.
	AsMap() Map
	// AsSlice returns a selector that selects Slices from the given
	// values.
	AsSlice() Slice
	// Number returns a selector taht selects Numbers from the given values.
	AsNumber() Number
	// String returns a selector that selects strings from the given values.
	AsString() String

	// Field returns the interface pointed by this specific field in the
	// map. If the field doesn't exist, the value will be filtered
	// out.
	Field(string) Interface
	// FieldP returns all the interfaces pointed by field that match the
	// string predicate. This selector can return more values than
	// it gets (for one map, it can returns multiple sub-values, one
	// for each field that matches the predicate).
	FieldP(...p.String) Interface

	// At returns a selector that select the child at the given
	// index, if the list has such an index. Otherwise, nothing is
	// returned.
	At(index int) Interface
	// AtP returns a selector that selects all the item whose index
	// matches the number predicate. More predicates can be given,
	// they are "and"-ed by this method.
	AtP(ips ...p.Number) Interface
	// Last returns a selector that selects the last value of the
	// list. If the list is empty, then nothing will be selected.
	Last() Interface

	// Children returns a selector that selects the direct children
	// of the given values.
	Children() Interface
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() Interface

	// Filter will create a new String that filters only the values
	// who match the predicate.
	Filter(...p.Interface) Interface
}

// Field returns the interface pointed by this specific field in the
// map. If the field doesn't exist, the value will be filtered
// out.
func Field(field string) Interface {
	return FieldP(p.StringEqual(field))
}

// FieldP returns all the interfaces pointed by field that match the
// string predicate. This selector can return more values than
// it gets (for one map, it can returns multiple sub-values, one
// for each field that matches the predicate).
func FieldP(predicates ...p.String) Interface {
	return &interfaceS{vf: interfaceFieldPFilter{sp: p.StringAnd(predicates...)}}
}

// At returns a selector that select the child at the given
// index, if the list has such an index. Otherwise, nothing is
// returned.
func At(index int) Interface {
	return AtP(p.NumberEqual(float64(index)))
}

// AtP returns a selector that selects all the item whose index
// matches the number predicate. More predicates can be given,
// they are "and"-ed by this method.
func AtP(ips ...p.Number) Interface {
	return &interfaceS{vf: interfaceAtPFilter{ip: p.NumberAnd(ips...)}}
}

// Last returns a selector that selects the last value of the
// list. If the list is empty, then nothing will be selected.
func Last() Interface {
	return &interfaceS{vf: interfaceLastFilter{}}
}

// Children selects all the children of the values.
func Children() Interface {
	return &interfaceS{vf: interfaceChildrenFilter{}}
}

// All selects all the direct and indirect childrens of the values.
func All() Interface {
	return &interfaceS{vf: interfaceAllFilter{}}
}

// Filter will only return the values that match the predicate.
func Filter(predicates ...p.Interface) Interface {
	return &interfaceS{vf: &interfaceFilterP{vp: p.InterfaceAnd(predicates...)}}
}

// Interface is a "Interface Selector". It selects a list of values, maps,
// slices, strings, integer from a list of values.
type interfaceS struct {
	vs Interface
	vf interfaceFilter
}

// Match returns true if the selector can find items in the given
// value. Otherwise, it returns false.
func (s *interfaceS) Match(value interface{}) bool {
	return len(s.SelectFrom(value)) != 0
}

func (s *interfaceS) SelectFrom(interfaces ...interface{}) []interface{} {
	if s.vs != nil {
		interfaces = s.vs.SelectFrom(interfaces...)
	}
	return s.vf.SelectFrom(interfaces...)
}

func (s *interfaceS) AsMap() Map {
	return &mapS{vs: s}
}

func (s *interfaceS) AsSlice() Slice {
	return &sliceS{vs: s}
}

func (s *interfaceS) AsNumber() Number {
	return &numberS{vs: s}
}

func (s *interfaceS) AsString() String {
	return &stringS{vs: s}
}

func (s *interfaceS) At(index int) Interface {
	return s.AtP(p.NumberEqual(float64(index)))
}

func (s *interfaceS) AtP(predicates ...p.Number) Interface {
	return &interfaceS{vs: s, vf: interfaceAtPFilter{ip: p.NumberAnd(predicates...)}}
}

func (s *interfaceS) Last() Interface {
	return &interfaceS{vs: s, vf: interfaceLastFilter{}}
}

func (s *interfaceS) Field(key string) Interface {
	return s.FieldP(p.StringEqual(key))
}

func (s *interfaceS) FieldP(predicates ...p.String) Interface {
	return &interfaceS{vs: s, vf: interfaceFieldPFilter{sp: p.StringAnd(predicates...)}}
}

func (s *interfaceS) Children() Interface {
	return &interfaceS{vs: s, vf: interfaceChildrenFilter{}}
}

func (s *interfaceS) All() Interface {
	return &interfaceS{vs: s, vf: interfaceAllFilter{}}
}

func (s *interfaceS) Filter(predicates ...p.Interface) Interface {
	return &interfaceS{vs: s, vf: &interfaceFilterP{vp: p.InterfaceAnd(predicates...)}}
}
