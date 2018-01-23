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

// NumberS is a "number selector". It selects values as numbers (if
// possible) and filters those numbers based on the "filtered"
// predicates.
type NumberS interface {
	// NumberS can be used as a Value predicate. If the selector can't
	// select any number from the value, then the predicate is
	// false.
	ValueP

	// SelectFrom finds numbers from values using this selector. The
	// list can be bigger or smaller than the initial lists,
	// depending on the select criterias.
	SelectFrom(...interface{}) []float64

	// Filter will create a new NumberS that filters only the values
	// who match the predicate.
	Filter(...NumberP) NumberS
}

// Number returns a NumberS that selects numbers from given values.
func Number() NumberS {
	return &numberS{}
}

type numberS struct {
	vs ValueS
	ip NumberP
}

func (s *numberS) SelectFrom(values ...interface{}) []float64 {
	numbers := []float64{}
	if s.vs != nil {
		values = s.vs.SelectFrom(values...)
	}
	for _, value := range values {
		i, ok := value.(float64)
		if !ok {
			continue
		}
		if s.ip != nil && !s.ip.Match(i) {
			continue
		}
		numbers = append(numbers, i)
	}
	return numbers
}

func (s *numberS) Filter(predicates ...NumberP) NumberS {
	if s.ip != nil {
		predicates = append(predicates, s.ip)
	}
	return &numberS{
		vs: s.vs,
		ip: NumberAnd(predicates...),
	}
}

func (s *numberS) Match(values interface{}) bool {
	return len(s.SelectFrom(values)) != 0
}
