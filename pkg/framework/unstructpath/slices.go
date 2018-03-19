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

// SliceS is a "slice selector". It selects values as slices (if
// possible) and filters those slices based on the "filtered"
// predicates.
type SliceS interface {
	// SliceS can be used as a Interface predicate. If the selector
	// can't select any slice from the value, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds slices from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) [][]interface{}

	// At returns a selector that select the child at the given
	// index, if the list has such an index. Otherwise, nothing is
	// returned.
	At(index int) InterfaceS
	// AtP returns a selector that selects all the item whose index
	// matches the number predicate. More predicates can be given,
	// they are "and"-ed by this method.
	AtP(ips ...p.Number) InterfaceS
	// Last returns a selector that selects the last value of the
	// list. If the list is empty, then nothing will be selected.
	Last() InterfaceS

	// All returns a selector that selects all direct and indrect
	// children of the given values.
	Children() InterfaceS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() InterfaceS

	// Filter will create a new SliceS that filters only the values
	// who match the predicate.
	Filter(...p.Slice) SliceS
}

// Slice creates a selector that takes values and filters them into
// slices if possible.
func Slice() SliceS {
	return &sliceS{}
}

type sliceS struct {
	vs InterfaceS
	sp p.Slice
}

func (s *sliceS) SelectFrom(interfaces ...interface{}) [][]interface{} {
	if s.vs != nil {
		interfaces = s.vs.SelectFrom(interfaces...)
	}

	slices := [][]interface{}{}
	for _, value := range interfaces {
		slice, ok := value.([]interface{})
		if !ok {
			continue
		}
		if s.sp != nil && !s.sp.Match(slice) {
			continue
		}
		slices = append(slices, slice)
	}

	return slices
}

func (s *sliceS) At(index int) InterfaceS {
	return s.AtP(p.NumberEqual(float64(index)))
}

func (s *sliceS) AtP(predicates ...p.Number) InterfaceS {
	return filterSlice(s, sliceAtPFilter{ip: p.NumberAnd(predicates...)})
}

func (s *sliceS) Last() InterfaceS {
	return filterSlice(s, sliceLastFilter{})
}

func (s *sliceS) Children() InterfaceS {
	// No predicates means select all direct children.
	return s.AtP()
}

func (s *sliceS) All() InterfaceS {
	return filterSlice(s, sliceAllFilter{})
}

func (s *sliceS) Filter(sps ...p.Slice) SliceS {
	return &sliceS{vs: s.vs, sp: p.SliceAnd(append(sps, s.sp)...)}
}

func (s *sliceS) Match(value interface{}) bool {
	return len(s.SelectFrom(value)) != 0
}
