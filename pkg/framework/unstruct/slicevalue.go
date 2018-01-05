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

package unstruct

type sliceValue struct {
	parent Value
	index  int
}

var _ Value = &sliceValue{}

func (s *sliceValue) Data() interface{} {
	return s.parent.Slice().Data()[s.index]
}

func (s *sliceValue) Parent() Value {
	return s.parent
}

func (s *sliceValue) Set(value interface{}) Value {
	s.parent.Slice().Data()[s.index] = value
	return s
}

func (s *sliceValue) Map() Map {
	return newMap(s)
}

func (s *sliceValue) Slice() Slice {
	return newSlice(s)
}
