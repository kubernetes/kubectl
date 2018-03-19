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

// MapP is a "map predicate". It's a type that decides if a
// map matches or not.
type MapP interface {
	Match(unstruct.Map) bool
}

// MapNot inverses the value of the predicate.
func MapNot(predicate MapP) MapP {
	return mapNot{mp: predicate}
}

type mapNot struct {
	mp MapP
}

func (p mapNot) Match(v unstruct.Map) bool {
	return !p.mp.Match(v)
}

// MapAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func MapAnd(predicates ...MapP) MapP {
	return mapAnd{mps: predicates}
}

type mapAnd struct {
	mps []MapP
}

func (p mapAnd) Match(m unstruct.Map) bool {
	for _, mp := range p.mps {
		if !mp.Match(m) {
			return false
		}
	}
	return true
}

// MapOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func MapOr(predicates ...MapP) MapP {
	mps := []MapP{}

	// Implements "De Morgan's law"
	for _, mp := range predicates {
		mps = append(mps, MapNot(mp))
	}
	return MapNot(MapAnd(mps...))
}

// MapNumFields matches if the number of fields matches the number
// predicate.
func MapNumFields(predicate NumberP) MapP {
	return mapNumFields{ip: predicate}
}

type mapNumFields struct {
	ip NumberP
}

func (p mapNumFields) Match(m unstruct.Map) bool {
	return p.ip.Match(float64(len(m.Data())))
}
