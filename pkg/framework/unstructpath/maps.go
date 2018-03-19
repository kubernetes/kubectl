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

// MapS is a "map selector". It selects interfaces as maps (if
// possible) and filters those maps based on the "filtered"
// predicates.
type MapS interface {
	// MapS can be used as a Interface predicate. If the selector can't
	// select any map from the interface, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds maps from interfaces using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) []map[string]interface{}

	// Field returns the interface pointed by this specific field in the
	// map. If the field doesn't exist, the value will be filtered
	// out.
	Field(string) InterfaceS
	// FieldP returns all the interfaces pointed by field that match the
	// string predicate. This selector can return more values than
	// it gets (for one map, it can returns multiple sub-values, one
	// for each field that matches the predicate).
	FieldP(...p.String) InterfaceS

	// All returns a selector that selects all direct and indrect
	// children of the given values.
	Children() InterfaceS
	// All returns a selector that selects all direct and indrect
	// children of the given values.
	All() InterfaceS

	// Filter will create a new MapS that filters only the values
	// who match the predicate.
	Filter(...p.Map) MapS
}

// Map creates a selector that takes interfaces and filters them into maps
// if possible.
func Map() MapS {
	return &mapS{}
}

type mapS struct {
	vs InterfaceS
	mp p.Map
}

func (s *mapS) SelectFrom(interfaces ...interface{}) []map[string]interface{} {
	if s.vs != nil {
		interfaces = s.vs.SelectFrom(interfaces...)
	}

	maps := []map[string]interface{}{}
	for _, value := range interfaces {
		m, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		if s.mp != nil && !s.mp.Match(m) {
			continue
		}
		maps = append(maps, m)
	}

	return maps
}

func (s *mapS) Field(str string) InterfaceS {
	return s.FieldP(p.StringEqual(str))
}

func (s *mapS) FieldP(predicates ...p.String) InterfaceS {
	return filterMap(s, mapFieldPFilter{sp: p.StringAnd(predicates...)})
}

func (s *mapS) Children() InterfaceS {
	// No predicate means select all.
	return s.FieldP()
}

func (s *mapS) All() InterfaceS {
	return filterMap(s, mapAllFilter{})
}

func (s *mapS) Filter(predicates ...p.Map) MapS {
	return &mapS{vs: s.vs, mp: p.MapAnd(append(predicates, s.mp)...)}
}

func (s *mapS) Match(value interface{}) bool {
	return len(s.SelectFrom(value)) != 0
}
