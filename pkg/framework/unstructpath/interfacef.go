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
	p "k8s.io/kubectl/pkg/framework/predicates"
)

// A interfaceFilter allows us to chain InterfaceS to InterfaceS. None of this is
// public. It's implementing the "SelectFrom" part of a InterfaceS.
type interfaceFilter interface {
	SelectFrom(...interface{}) []interface{}
}

// interfaceFilterP filters using a predicate.
type interfaceFilterP struct {
	vp p.Interface
}

func (f *interfaceFilterP) SelectFrom(interfaces ...interface{}) []interface{} {
	vs := []interface{}{}
	for _, value := range interfaces {
		if f.vp.Match(value) {
			vs = append(vs, value)
		}
	}
	return vs
}

type interfaceChildrenFilter struct{}

func (interfaceChildrenFilter) SelectFrom(interfaces ...interface{}) []interface{} {
	children := []interface{}{}
	// We could process all slices and then all maps, but we want to
	// keep things in the same order.
	for _, value := range interfaces {
		// Only one of the two should do something useful.
		children = append(children, Slice().Children().SelectFrom(value)...)
		children = append(children, Map().Children().SelectFrom(value)...)
	}
	return children
}

// interfaceSliceFilter is a Interface-to-Slice combined with a Slice-to-Interface
// to form a Interface-to-Interface.
type interfaceSliceFilter struct {
	ss SliceS
	sf sliceFilter
}

func (s *interfaceSliceFilter) SelectFrom(interfaces ...interface{}) []interface{} {
	return s.sf.SelectFrom(s.ss.SelectFrom(interfaces...)...)
}

// interfaceMapFilter is a Interface-to-Map combined with a Map-to-Interface to form
// a Interface-to-Interface.
type interfaceMapFilter struct {
	ms MapS
	mf mapFilter
}

func (s *interfaceMapFilter) SelectFrom(interfaces ...interface{}) []interface{} {
	return s.mf.SelectFrom(s.ms.SelectFrom(interfaces...)...)
}

type interfaceAllFilter struct{}

func (interfaceAllFilter) SelectFrom(interfaces ...interface{}) []interface{} {
	vs := []interface{}{}
	for _, value := range interfaces {
		vs = append(vs, value)
		// Only one of the follow two statements should return something ...
		vs = append(vs, Slice().All().SelectFrom(value)...)
		vs = append(vs, Map().All().SelectFrom(value)...)
	}
	return vs
}
