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

// StringS is a "string selector". It selects values as strings (if
// possible) and filters those strings based on the "filtered"
// predicates.
type StringS interface {
	// StringS can be used as a Interface predicate. If the selector can't
	// select any string from the value, then the predicate is
	// false.
	p.Interface

	// SelectFrom finds strings from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) []string

	// Filter will create a new StringS that filters only the values
	// who match the predicate.
	Filter(...p.String) StringS
}

type stringS struct {
	vs InterfaceS
	sp p.String
}

// String returns a StringS that selects strings from values.
func String() StringS {
	return &stringS{}
}

func (s *stringS) SelectFrom(values ...interface{}) []string {
	strings := []string{}
	if s.vs != nil {
		values = s.vs.SelectFrom(values...)
	}
	for _, value := range values {
		str, ok := value.(string)
		if !ok {
			continue
		}
		if s.sp != nil && !s.sp.Match(str) {
			continue
		}
		strings = append(strings, str)
	}
	return strings
}

func (s *stringS) Filter(predicates ...p.String) StringS {
	if s.sp != nil {
		predicates = append(predicates, s.sp)
	}
	return &stringS{
		vs: s.vs,
		sp: p.StringAnd(predicates...),
	}
}

func (s *stringS) Match(values interface{}) bool {
	return len(s.SelectFrom(values)) != 0
}
