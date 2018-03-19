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
	"k8s.io/kubectl/pkg/framework/unstruct"
)

// SliceS is a "slice selector". It selects values as slices (if
// possible) and filters those slices based on the "filtered"
// predicates.
type SliceS interface {
	// SliceS can be used as a Value predicate. If the selector
	// can't select any slice from the value, then the predicate is
	// false.
	ValueP

	// Select finds slices from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	Select(...unstruct.Value) []unstruct.Slice

	// At returns a selector that select the child at the given
	// index, if the list has such an index. Otherwise, nothing is
	// returned.
	At(index int) ValueS
	// AtP returns a selector that selects all the item whose index
	// matches the number predicate. More predicates can be given,
	// they are "and"-ed by this method.
	AtP(ips ...NumberP) ValueS
	// Last returns a selector that selects the last value of the
	// list. If the list is empty, then nothing will be selected.
	Last() ValueS

	// Parent returns a selector that selects the parent of each
	// given values. If the value is a root, then no value is
	// selected.
	Parent() ValueS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	Children() ValueS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() ValueS

	// Filter will create a new SliceS that filters only the values
	// who match the predicate.
	Filter(...SliceP) SliceS
}

// Slice creates a selector that takes values and filters them into
// slices if possible.
func Slice() SliceS {
	return &sliceS{}
}

type sliceS struct {
	vs ValueS
	sp SliceP
}

func (s *sliceS) Select(values ...unstruct.Value) []unstruct.Slice {
	if s.vs != nil {
		values = s.vs.Select(values...)
	}

	slices := []unstruct.Slice{}
	for _, value := range values {
		slice := value.Slice()
		if slice == nil {
			continue
		}
		if s.sp != nil && !s.sp.Match(slice) {
			continue
		}
		slices = append(slices, slice)
	}

	return slices
}

func (s *sliceS) At(index int) ValueS {
	return s.AtP(NumberEqual(float64(index)))
}

func (s *sliceS) AtP(predicates ...NumberP) ValueS {
	return filterSlice(s, sliceAtPFilter{ip: NumberAnd(predicates...)})
}

func (s *sliceS) Last() ValueS {
	return filterSlice(s, sliceLastFilter{})
}

func (s *sliceS) Parent() ValueS {
	return filterSlice(s, sliceParentFilter{})
}

func (s *sliceS) Children() ValueS {
	// No predicates means select all direct children.
	return s.AtP()
}

func (s *sliceS) All() ValueS {
	return filterSlice(s, sliceAllFilter{})
}

func (s *sliceS) Filter(sps ...SliceP) SliceS {
	return &sliceS{vs: s.vs, sp: SliceAnd(append(sps, s.sp)...)}
}

func (s *sliceS) Match(value unstruct.Value) bool {
	return len(s.Select(value)) != 0
}
