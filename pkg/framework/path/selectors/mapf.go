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

// This is a Map-to-Interface filter.
type mapFilter interface {
	SelectFrom(...map[string]interface{}) []interface{}
}

func filterMap(ms MapS, mf mapFilter) InterfaceS {
	return &interfaceS{
		vf: &interfaceMapFilter{
			ms: ms,
			mf: mf,
		},
	}
}

type mapFieldPFilter struct {
	sp p.String
}

func (f mapFieldPFilter) SelectFrom(maps ...map[string]interface{}) []interface{} {
	interfaces := []interface{}{}

	for _, m := range maps {
		for _, field := range sortedKeys(m) {
			if !f.sp.Match(field) {
				continue
			}
			interfaces = append(interfaces, m[field])
		}
	}
	return interfaces
}

type mapAllFilter struct{}

func (mapAllFilter) SelectFrom(maps ...map[string]interface{}) []interface{} {
	interfaces := []interface{}{}
	for _, m := range maps {
		for _, field := range sortedKeys(m) {
			interfaces = append(interfaces, All().SelectFrom(m[field])...)
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
