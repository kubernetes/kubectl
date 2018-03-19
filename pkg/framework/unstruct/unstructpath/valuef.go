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

// A valueFilter allows us to chain ValueS to ValueS. None of this is
// public. It's implementing the "Select" part of a ValueS.
type valueFilter interface {
	Select(...unstruct.Value) []unstruct.Value
}

// valueFilterP filters using a predicate.
type valueFilterP struct {
	vp ValueP
}

func (f *valueFilterP) Select(values ...unstruct.Value) []unstruct.Value {
	vs := []unstruct.Value{}
	for _, value := range values {
		if f.vp.Match(value) {
			vs = append(vs, value)
		}
	}
	return vs
}

type valueParentFilter struct{}

func (valueParentFilter) Select(values ...unstruct.Value) []unstruct.Value {
	vs := []unstruct.Value{}
	for _, value := range values {
		if p := value.Parent(); p != nil {
			vs = append(vs, p)
		}
	}
	return vs
}

type valueRootFilter struct{}

func (valueRootFilter) Select(values ...unstruct.Value) []unstruct.Value {
	vs := []unstruct.Value{}
	for _, value := range values {
		root := unstruct.Root(value)
		if root == nil {
			continue
		}
		vs = append(vs, root)
	}
	return vs
}

type valueChildrenFilter struct{}

func (valueChildrenFilter) Select(values ...unstruct.Value) []unstruct.Value {
	children := []unstruct.Value{}
	// We could process all slices and then all maps, but we want to
	// keep things in the same order.
	for _, value := range values {
		// Only one of the two should do something useful.
		children = append(children, Slice().Children().Select(value)...)
		children = append(children, Map().Children().Select(value)...)
	}
	return children
}

// valueSliceFilter is a Value-to-Slice combined with a Slice-to-Value
// to form a Value-to-Value.
type valueSliceFilter struct {
	ss SliceS
	sf sliceFilter
}

func (s *valueSliceFilter) Select(values ...unstruct.Value) []unstruct.Value {
	return s.sf.Select(s.ss.Select(values...)...)
}

// valueMapFilter is a Value-to-Map combined with a Map-to-Value to form
// a Value-to-Value.
type valueMapFilter struct {
	ms MapS
	mf mapFilter
}

func (s *valueMapFilter) Select(values ...unstruct.Value) []unstruct.Value {
	return s.mf.Select(s.ms.Select(values...)...)
}

type valueAllFilter struct{}

func (valueAllFilter) Select(values ...unstruct.Value) []unstruct.Value {
	vs := []unstruct.Value{}
	for _, value := range values {
		vs = append(vs, value)
		// Only one of the follow two statements should return something ...
		vs = append(vs, Slice().All().Select(value)...)
		vs = append(vs, Map().All().Select(value)...)
	}
	return vs
}
