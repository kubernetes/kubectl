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
	"fmt"
)

// A new filter requires a type with the two methods required by the
// Filter interface.

type filterLetterP struct {
}

func (*filterLetterP) Resource(r *Resource) bool {
	if r.Resource.Name != "" {
		return string(r.Resource.Name[0]) == "p"
	} else {
		return false
	}
}

func (*filterLetterP) SubResource(*SubResource) bool {
	return true
}

// Example_filter demonstrates how to create and then apply a resource
// filter. This (admittedly ludicrous) example should have the same
// output as the Resources example, but it excludes all resources that
// do not start with the letter "p".
func Example_filter() {

	// This is a blank parser for testing only. Use actual values in
	// your code.
	p := NewParser(nil, nil, "", "")
	r, err := p.Resources()
	if err != nil {
		panic(err)
	}
	// Applying the filter...
	r = r.Filter(&filterLetterP{})
	for name, versions := range r {
		fmt.Println("\n→", name)
		for _, version := range versions {
			fmt.Printf("→→ %+v\n", version)
		}
	}
}
