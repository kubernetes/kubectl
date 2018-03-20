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
	p "k8s.io/kubectl/pkg/framework/path/predicates"
)

// Map is a "map selector". It selects interfaces as maps (if
// possible) and filters those maps based on the "filtered"
// predicates.
type Map interface {
	// Map can be used as a Interface predicate. If the selector can't
	// select any map from the interface, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds maps from interfaces using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) []map[string]interface{}

	// Filter will create a new Map that filters only the values
	// who match the predicate.
	Filter(...p.Map) Map
}

// Map creates a selector that takes interfaces and filters them into maps
// if possible.
func AsMap() Map {
	return &mapS{}
}

type mapS struct {
	vs Interface
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

func (s *mapS) Filter(predicates ...p.Map) Map {
	return &mapS{vs: s.vs, mp: p.MapAnd(append(predicates, s.mp)...)}
}

func (s *mapS) Match(value interface{}) bool {
	return len(s.SelectFrom(value)) != 0
}
