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

import (
	"reflect"
)

// Interface is a "interface{} predicate". It's a type that decides if an
// interface matches or not.
type Interface interface {
	Match(interface{}) bool
}

// InterfaceDeepEqual compares the Interface data with DeepEqual.
func InterfaceDeepEqual(v interface{}) Interface {
	return interfaceEqual{v: v}
}

type interfaceEqual struct {
	v interface{}
}

func (p interfaceEqual) Match(v interface{}) bool {
	return reflect.DeepEqual(v, p.v)
}

// InterfaceNot inverses the value of the predicate.
func InterfaceNot(predicate Interface) Interface {
	return interfaceNot{vp: predicate}
}

type interfaceNot struct {
	vp Interface
}

func (p interfaceNot) Match(v interface{}) bool {
	return !p.vp.Match(v)
}

// InterfaceAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func InterfaceAnd(predicates ...Interface) Interface {
	return interfaceAnd{vps: predicates}
}

type interfaceAnd struct {
	vps []Interface
}

func (p interfaceAnd) Match(i interface{}) bool {
	for _, vp := range p.vps {
		if !vp.Match(i) {
			return false
		}
	}
	return true
}

// InterfaceOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func InterfaceOr(predicates ...Interface) Interface {
	vps := []Interface{}

	// Implements "De Morgan's law"
	for _, vp := range predicates {
		vps = append(vps, InterfaceNot(vp))
	}
	return InterfaceNot(InterfaceAnd(vps...))
}
