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

// Map is a "map predicate". It's a type that decides if a
// map matches or not.
type Map interface {
	Match(map[string]interface{}) bool
}

// MapNot inverses the value of the predicate.
func MapNot(predicate Map) Map {
	return mapNot{mp: predicate}
}

type mapNot struct {
	mp Map
}

func (p mapNot) Match(m map[string]interface{}) bool {
	return !p.mp.Match(m)
}

// MapAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func MapAnd(predicates ...Map) Map {
	return mapAnd{mps: predicates}
}

type mapAnd struct {
	mps []Map
}

func (p mapAnd) Match(m map[string]interface{}) bool {
	for _, mp := range p.mps {
		if !mp.Match(m) {
			return false
		}
	}
	return true
}

// MapOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func MapOr(predicates ...Map) Map {
	mps := []Map{}

	// Implements "De Morgan's law"
	for _, mp := range predicates {
		mps = append(mps, MapNot(mp))
	}
	return MapNot(MapAnd(mps...))
}

// MapNumFields matches if the number of fields matches the number
// predicate.
func MapNumFields(predicate Number) Map {
	return mapNumFields{ip: predicate}
}

type mapNumFields struct {
	ip Number
}

func (p mapNumFields) Match(m map[string]interface{}) bool {
	return p.ip.Match(float64(len(m)))
}
