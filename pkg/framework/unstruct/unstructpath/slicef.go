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
	Select(...unstruct.Slice) []unstruct.Value
}

type sliceAtPFilter struct {
	ip NumberP
}

func (f sliceAtPFilter) Select(slices ...unstruct.Slice) []unstruct.Value {
	values := []unstruct.Value{}

	for _, slice := range slices {
		for i := 0; i < slice.Length(); i++ {
			if !f.ip.Match(float64(i)) {
				continue
			}
			values = append(values, slice.At(i))
		}
	}
	return values
}

type sliceParentFilter struct{}

func (f sliceParentFilter) Select(slices ...unstruct.Slice) []unstruct.Value {
	values := []unstruct.Value{}
	for _, slice := range slices {
		if p := slice.Parent(); p != nil {
			values = append(values, p)
		}
	}
	return values
}

type sliceLastFilter struct{}

func (f sliceLastFilter) Select(slices ...unstruct.Slice) []unstruct.Value {
	values := []unstruct.Value{}
	for _, slice := range slices {
		if slice.Length() == 0 {
			continue
		}
		values = append(values, slice.At(slice.Length()-1))
	}
	return values
}

type sliceAllFilter struct{}

func (sliceAllFilter) Select(slices ...unstruct.Slice) []unstruct.Value {
	values := []unstruct.Value{}
	for _, slice := range slices {
		for i := 0; i < slice.Length(); i++ {
			values = append(values, All().Select(slice.At(i))...)
		}
	}
	return values
}
