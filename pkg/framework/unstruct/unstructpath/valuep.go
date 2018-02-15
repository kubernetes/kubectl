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
	"reflect"

	"k8s.io/kubectl/pkg/framework/unstruct"
)

// ValueP is a "value predicate". It's a type that decides if a
// value matches or not.
type ValueP interface {
	Match(unstruct.Value) bool
}

// ValueDeepEqual compares the Value data with DeepEqual.
func ValueDeepEqual(v unstruct.Value) ValueP {
	return valueEqual{v: v}
}

type valueEqual struct {
	v unstruct.Value
}

func (p valueEqual) Match(v unstruct.Value) bool {
	return reflect.DeepEqual(v.Data(), p.v.Data())
}

// ValueNot inverses the value of the predicate.
func ValueNot(predicate ValueP) ValueP {
	return valueNot{vp: predicate}
}

type valueNot struct {
	vp ValueP
}

func (p valueNot) Match(v unstruct.Value) bool {
	return !p.vp.Match(v)
}

// ValueAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func ValueAnd(predicates ...ValueP) ValueP {
	return valueAnd{vps: predicates}
}

type valueAnd struct {
	vps []ValueP
}

func (p valueAnd) Match(value unstruct.Value) bool {
	for _, vp := range p.vps {
		if !vp.Match(value) {
			return false
		}
	}
	return true
}

// ValueOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func ValueOr(predicates ...ValueP) ValueP {
	vps := []ValueP{}

	// Implements "De Morgan's law"
	for _, vp := range predicates {
		vps = append(vps, ValueNot(vp))
	}
	return ValueNot(ValueAnd(vps...))
}
