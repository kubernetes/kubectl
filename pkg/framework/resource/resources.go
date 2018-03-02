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

import (
	"sort"
)

// Resources is the set of resources found in the API server.
type Resources map[string][]*Resource

// Names returns the resource names ordered alphabetically.
func (r Resources) Names() []string {
	sorted := []string{}
	for k := range r {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

// Filter takes a Filter interface which is used to filter the resources
// and/or subresources in a Resources object. It returns a new Resources
// object.
func (r Resources) Filter(filter Filter) Resources {
	filtered := Resources{}
	for k, v := range r {
		for _, version := range v {
			if !filter.Resource(version) {
				continue
			}
			copy := r.filterSubResources(*version, filter)
			filtered[k] = append(filtered[k], &copy)
		}
	}
	return filtered
}

// filterSubresources returns a copy of resource with the subresources
// filtered.
func (r Resources) filterSubResources(resource Resource, filter Filter) Resource {
	filtered := []*SubResource{}
	for _, v := range resource.SubResources {
		if !filter.SubResource(v) {
			continue
		}
		filtered = append(filtered, v)
	}
	resource.SubResources = filtered
	return resource
}
