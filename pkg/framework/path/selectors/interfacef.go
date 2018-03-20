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
	"sort"

	p "k8s.io/kubectl/pkg/framework/path/predicates"
)

// A interfaceFilter allows us to chain Interface to Interface. None of this is
// public. It's implementing the "SelectFrom" part of a Interface.
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
		// AtP() with nothing selects all the children of a list
		children = append(children, AtP().SelectFrom(value)...)
		// FieldP() with nothing selects all the children of a map
		children = append(children, FieldP().SelectFrom(value)...)
	}
	return children
}

type interfaceAtPFilter struct {
	ip p.Number
}

func (f interfaceAtPFilter) SelectFrom(values ...interface{}) []interface{} {
	interfaces := []interface{}{}

	for _, value := range values {
		slice, ok := value.([]interface{})
		if !ok {
			continue
		}
		for i := range slice {
			if !f.ip.Match(float64(i)) {
				continue
			}
			interfaces = append(interfaces, slice[i])
		}
	}
	return interfaces
}

type interfaceLastFilter struct{}

func (f interfaceLastFilter) SelectFrom(values ...interface{}) []interface{} {
	interfaces := []interface{}{}
	for _, value := range values {
		slice, ok := value.([]interface{})
		if !ok {
			continue
		}
		if len(slice) == 0 {
			continue
		}
		interfaces = append(interfaces, slice[len(slice)-1])
	}
	return interfaces
}

type interfaceFieldPFilter struct {
	sp p.String
}

func (f interfaceFieldPFilter) SelectFrom(values ...interface{}) []interface{} {
	interfaces := []interface{}{}

	for _, value := range values {
		m, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		for _, field := range sortedKeys(m) {
			if !f.sp.Match(field) {
				continue
			}
			interfaces = append(interfaces, m[field])
		}
	}
	return interfaces
}

func sortedKeys(m map[string]interface{}) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

type interfaceAllFilter struct{}

func (interfaceAllFilter) SelectFrom(interfaces ...interface{}) []interface{} {
	vs := []interface{}{}
	for _, value := range interfaces {
		vs = append(vs, value)
		vs = append(vs, Children().All().SelectFrom(value)...)
	}
	return vs
}
