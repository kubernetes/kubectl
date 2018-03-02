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

package resource

// Filter is an interface whose methods are used to arbitrarily filter
// resources and subresources in a Resources object. Filtering criteria
// can be anything at all, so long as the "accepted" resources return
// true, and the filtered out resources return false.
type Filter interface {
	Resource(*Resource) bool
	SubResource(*SubResource) bool
}

// NewAndFilter that takes as argument a slice of filters, which then
// work as a logical AND. This filter will only return resources that
// pass through every filter in the argument slice. It returns a new
// Filter value.
func NewAndFilter(filters ...Filter) Filter {
	return &andFilter{Filters: filters}
}

// NewOrFilter that takes as argument a slice of filters, which then work
// as a logical OR. This filter will return resources that pass through
// any one or more filters in the argument slice. It returns a new Filter
// value.
func NewOrFilter(filters ...Filter) Filter {
	return &orFilter{Filters: filters}
}

type emptyFilter struct {
}

func (*emptyFilter) Resource(*Resource) bool {
	return true
}

func (*emptyFilter) SubResource(*SubResource) bool {
	return true
}

type andFilter struct {
	Filters []Filter
}

func (a *andFilter) Resource(r *Resource) bool {
	for _, f := range a.Filters {
		if !f.Resource(r) {
			return false
		}
	}
	return true
}

func (a *andFilter) SubResource(sr *SubResource) bool {
	for _, f := range a.Filters {
		if !f.SubResource(sr) {
			return false
		}
	}
	return true
}

type orFilter struct {
	Filters []Filter
}

func (a *orFilter) Resource(r *Resource) bool {
	for _, f := range a.Filters {
		if f.Resource(r) {
			return true
		}
	}
	return false
}

func (a *orFilter) SubResource(sr *SubResource) bool {
	for _, f := range a.Filters {
		if f.SubResource(sr) {
			return true
		}
	}
	return false
}
