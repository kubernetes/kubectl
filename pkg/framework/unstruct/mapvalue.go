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

type mapValue struct {
	parent Value
	key    string
}

var _ Value = &mapValue{}

func (b *mapValue) Data() interface{} {
	return b.parent.Map().Data()[b.key]
}

func (b *mapValue) Parent() Value {
	return b.parent
}

func (b *mapValue) Set(value interface{}) Value {
	b.parent.Map().Data()[b.key] = value
	return b
}

func (b *mapValue) Map() Map {
	return newMap(b)
}

func (b *mapValue) Slice() Slice {
	return newSlice(b)
}
