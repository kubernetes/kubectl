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

package predicates

// Slice is a "slice predicate". It's a type that decides if a
// slice matches or not.
type Slice interface {
	Match([]interface{}) bool
}

// SliceNot inverses the value of the predicate.
func SliceNot(predicate Slice) Slice {
	return sliceNot{vp: predicate}
}

type sliceNot struct {
	vp Slice
}

func (p sliceNot) Match(slice []interface{}) bool {
	return !p.vp.Match(slice)
}

// SliceAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func SliceAnd(predicates ...Slice) Slice {
	return sliceAnd{sps: predicates}
}

type sliceAnd struct {
	sps []Slice
}

func (p sliceAnd) Match(slice []interface{}) bool {
	for _, sp := range p.sps {
		if !sp.Match(slice) {
			return false
		}
	}
	return true
}

// SliceOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func SliceOr(predicates ...Slice) Slice {
	sps := []Slice{}

	// Implements "De Morgan's law"
	for _, sp := range predicates {
		sps = append(sps, SliceNot(sp))
	}
	return SliceNot(SliceAnd(sps...))
}

// SliceLength matches if the length of the list matches the given
// integer predicate.
func SliceLength(ip Number) Slice {
	return sliceLength{ip: ip}
}

type sliceLength struct {
	ip Number
}

func (p sliceLength) Match(slice []interface{}) bool {
	return p.ip.Match(float64(len(slice)))
}
