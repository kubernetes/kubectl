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

import "sort"

type object struct {
	value Value
}

var _ Map = &object{}

func newMap(value Value) Map {
	if value.Data() == nil {
		value.Set(map[string]interface{}{})
	}
	_, ok := value.Data().(map[string]interface{})
	if !ok {
		return nil
	}
	return &object{value: value}
}

func (o *object) Data() map[string]interface{} {
	d, ok := o.value.Data().(map[string]interface{})
	if !ok {
		return nil
	}
	return d
}

func (o *object) Parent() Value {
	return o.value.Parent()
}

func (o *object) Field(key string) Value {
	if !o.HasField(key) {
		o.Data()[key] = nil
	}
	return &mapValue{parent: o.value, key: key}
}

func (o *object) HasField(key string) bool {
	_, ok := o.Data()[key]
	return ok
}

func (o *object) Fields() []string {
	fields := []string{}
	for key := range o.Data() {
		fields = append(fields, key)
	}
	sort.Strings(fields)
	return fields
}
