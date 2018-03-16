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

// SliceP is a "slice predicate". It's a type that decides if a
// slice matches or not.
type SliceP interface {
	Match([]interface{}) bool
}

// SliceNot inverses the value of the predicate.
func SliceNot(predicate SliceP) SliceP {
	return sliceNot{vp: predicate}
}

type sliceNot struct {
	vp SliceP
}

func (p sliceNot) Match(slice []interface{}) bool {
	return !p.vp.Match(slice)
}

// SliceAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func SliceAnd(predicates ...SliceP) SliceP {
	return sliceAnd{sps: predicates}
}

type sliceAnd struct {
	sps []SliceP
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
func SliceOr(predicates ...SliceP) SliceP {
	sps := []SliceP{}

	// Implements "De Morgan's law"
	for _, sp := range predicates {
		sps = append(sps, SliceNot(sp))
	}
	return SliceNot(SliceAnd(sps...))
}

// SliceLength matches if the length of the list matches the given
// integer predicate.
func SliceLength(ip NumberP) SliceP {
	return sliceLength{ip: ip}
}

type sliceLength struct {
	ip NumberP
}

func (p sliceLength) Match(slice []interface{}) bool {
	return p.ip.Match(float64(len(slice)))
}
