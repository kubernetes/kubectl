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

// Slice is a "slice selector". It selects values as slices (if
// possible) and filters those slices based on the "filtered"
// predicates.
type Slice interface {
	// Slice can be used as a Interface predicate. If the selector
	// can't select any slice from the value, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds slices from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) [][]interface{}

	// Filter will create a new Slice that filters only the values
	// who match the predicate.
	Filter(...p.Slice) Slice
}

// Slice creates a selector that takes values and filters them into
// slices if possible.
func AsSlice() Slice {
	return &sliceS{}
}

type sliceS struct {
	vs Interface
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

func (s *sliceS) Filter(sps ...p.Slice) Slice {
	return &sliceS{vs: s.vs, sp: p.SliceAnd(append(sps, s.sp)...)}
}

func (s *sliceS) Match(value interface{}) bool {
	return len(s.SelectFrom(value)) != 0
}
