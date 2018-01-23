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

import "k8s.io/kubectl/pkg/framework/unstruct"

// This is a Map-to-Value filter.
type mapFilter interface {
	Select(...unstruct.Map) []unstruct.Value
}

func filterMap(ms MapS, mf mapFilter) ValueS {
	return &valueS{
		vf: &valueMapFilter{
			ms: ms,
			mf: mf,
		},
	}
}

type mapFieldPFilter struct {
	sp StringP
}

func (f mapFieldPFilter) Select(maps ...unstruct.Map) []unstruct.Value {
	values := []unstruct.Value{}

	for _, m := range maps {
		for _, field := range m.Fields() {
			if !f.sp.Match(field) {
				continue
			}
			values = append(values, m.Field(field))
		}
	}
	return values
}

type mapParentFilter struct{}

func (mapParentFilter) Select(maps ...unstruct.Map) []unstruct.Value {
	values := []unstruct.Value{}
	for _, m := range maps {
		if p := m.Parent(); p != nil {
			values = append(values, p)
		}
	}
	return values
}

type mapAllFilter struct{}

func (mapAllFilter) Select(maps ...unstruct.Map) []unstruct.Value {
	values := []unstruct.Value{}
	for _, m := range maps {
		for _, field := range m.Fields() {
			values = append(values, All().Select(m.Field(field))...)
		}
	}
	return values
}
