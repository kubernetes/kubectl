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

func filterSlice(ss SliceS, sf sliceFilter) ValueS {
	return &valueS{
		vf: &valueSliceFilter{
			ss: ss,
			sf: sf,
		},
	}
}

// This is a Slice-to-Value filter.
type sliceFilter interface {
	SelectFrom(...[]interface{}) []interface{}
}

type sliceAtPFilter struct {
	ip p.Number
}

func (f sliceAtPFilter) SelectFrom(slices ...[]interface{}) []interface{} {
	values := []interface{}{}

	for _, slice := range slices {
		for i := range slice {
			if !f.ip.Match(float64(i)) {
				continue
			}
			values = append(values, slice[i])
		}
	}
	return values
}

type sliceLastFilter struct{}

func (f sliceLastFilter) SelectFrom(slices ...[]interface{}) []interface{} {
	values := []interface{}{}
	for _, slice := range slices {
		if len(slice) == 0 {
			continue
		}
		values = append(values, slice[len(slice)-1])
	}
	return values
}

type sliceAllFilter struct{}

func (sliceAllFilter) SelectFrom(slices ...[]interface{}) []interface{} {
	values := []interface{}{}
	for _, slice := range slices {
		for _, v := range slice {
			values = append(values, All().SelectFrom(v)...)
		}
	}
	return values
}
