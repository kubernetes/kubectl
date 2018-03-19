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

// InterfaceS is a "interface selector". It filters interfaces based on the
// "filtered" predicates.
type InterfaceS interface {
	// InterfaceS can be used as a Interface predicate. If the selector can't
	// select any interface from the interface, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds interfaces from interfaces using this selector. The
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
	Children() InterfaceS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() InterfaceS

	// Filter will create a new StringS that filters only the values
	// who match the predicate.
	Filter(...p.Interface) InterfaceS
}

// Children selects all the children of the values.
func Children() InterfaceS {
	return &interfaceS{vf: interfaceChildrenFilter{}}
}

// All selects all the direct and indirect childrens of the values.
func All() InterfaceS {
	return &interfaceS{vf: interfaceAllFilter{}}
}

// Filter will only return the values that match the predicate.
func Filter(predicates ...p.Interface) InterfaceS {
	return &interfaceS{vf: &interfaceFilterP{vp: p.InterfaceAnd(predicates...)}}
}

// InterfaceS is a "Interface SelectFromor". It selects a list of values, maps,
// slices, strings, integer from a list of values.
type interfaceS struct {
	vs InterfaceS
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

func (s *interfaceS) Map() MapS {
	return &mapS{vs: s}
}

func (s *interfaceS) Slice() SliceS {
	return &sliceS{vs: s}
}

func (s *interfaceS) Number() NumberS {
	return &numberS{vs: s}
}

func (s *interfaceS) String() StringS {
	return &stringS{vs: s}
}

func (s *interfaceS) Children() InterfaceS {
	return &interfaceS{vs: s, vf: interfaceChildrenFilter{}}
}

func (s *interfaceS) All() InterfaceS {
	return &interfaceS{vs: s, vf: interfaceAllFilter{}}
}

func (s *interfaceS) Filter(predicates ...p.Interface) InterfaceS {
	return &interfaceS{vs: s, vf: &interfaceFilterP{vp: p.InterfaceAnd(predicates...)}}
}
