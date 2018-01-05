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

type slice struct {
	value Value
}

var _ Slice = &slice{}

func newSlice(value Value) Slice {
	if value.Data() == nil {
		value.Set([]interface{}{})
	}
	_, ok := value.Data().([]interface{})
	if !ok {
		return nil
	}
	return &slice{value: value}
}

func (s *slice) Data() []interface{} {
	d, ok := s.value.Data().([]interface{})
	if !ok {
		return nil
	}
	return d
}

func (s *slice) Parent() Value {
	return s.value.Parent()
}

func (s *slice) Length() int {
	return len(s.Data())
}

func (s *slice) At(index int) Value {
	if index < 0 || index >= s.Length() {
		return nil
	}
	return &sliceValue{parent: s.value, index: index}
}

func (s *slice) Append(value interface{}) Value {
	s.value.Set(append(s.Data(), value))
	return s.At(s.Length() - 1)
}

func (s *slice) InsertAt(index int, value interface{}) Slice {
	if index < 0 || index >= s.Length() {
		return nil
	}
	a := s.Data()
	a = append(a, nil)
	copy(a[index+1:], a[index:])
	a[index] = value
	s.value.Set(a)
	return s

}
