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

// A valueFilter allows us to chain ValueS to ValueS. None of this is
// public. It's implementing the "SelectFrom" part of a ValueS.
type valueFilter interface {
	SelectFrom(...interface{}) []interface{}
}

// valueFilterP filters using a predicate.
type valueFilterP struct {
	vp ValueP
}

func (f *valueFilterP) SelectFrom(values ...interface{}) []interface{} {
	vs := []interface{}{}
	for _, value := range values {
		if f.vp.Match(value) {
			vs = append(vs, value)
		}
	}
	return vs
}

type valueChildrenFilter struct{}

func (valueChildrenFilter) SelectFrom(values ...interface{}) []interface{} {
	children := []interface{}{}
	// We could process all slices and then all maps, but we want to
	// keep things in the same order.
	for _, value := range values {
		// Only one of the two should do something useful.
		children = append(children, Slice().Children().SelectFrom(value)...)
		children = append(children, Map().Children().SelectFrom(value)...)
	}
	return children
}

// valueSliceFilter is a Value-to-Slice combined with a Slice-to-Value
// to form a Value-to-Value.
type valueSliceFilter struct {
	ss SliceS
	sf sliceFilter
}

func (s *valueSliceFilter) SelectFrom(values ...interface{}) []interface{} {
	return s.sf.SelectFrom(s.ss.SelectFrom(values...)...)
}

// valueMapFilter is a Value-to-Map combined with a Map-to-Value to form
// a Value-to-Value.
type valueMapFilter struct {
	ms MapS
	mf mapFilter
}

func (s *valueMapFilter) SelectFrom(values ...interface{}) []interface{} {
	return s.mf.SelectFrom(s.ms.SelectFrom(values...)...)
}

type valueAllFilter struct{}

func (valueAllFilter) SelectFrom(values ...interface{}) []interface{} {
	vs := []interface{}{}
	for _, value := range values {
		vs = append(vs, value)
		// Only one of the follow two statements should return something ...
		vs = append(vs, Slice().All().SelectFrom(value)...)
		vs = append(vs, Map().All().SelectFrom(value)...)
	}
	return vs
}
