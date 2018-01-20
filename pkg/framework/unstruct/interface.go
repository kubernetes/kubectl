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

// Map lets you manipulate an object.
type Map interface {
	// Data returns the underlying data provided by the object.
	Data() map[string]interface{}

	// Parent returns the parent value or nil if no parent is
	// present (root).
	Parent() Value

	// Field returns the Value object in the "key" field.
	//
	// If this key is not present in the map, this function will
	// create, set and return a new Value for this field, with a nil
	// content.
	Field(key string) Value

	// HasField returns true if the field exists.
	HasField(key string) bool

	// Fields returns the sorted list of fields in the map.
	Fields() []string
}

// Slice lets you manipulate an array.
type Slice interface {
	// Data returns the underlying data provided by the object.
	Data() []interface{}

	// Parent returns the parent value or nil if no parent is
	// present (root).
	Parent() Value

	// Length returns the number of items in the slice.
	Length() int

	// At returns the "index"-th item in the slice.
	//
	// If the index is out-of-bound, nil is returned.
	At(index int) Value

	// Append adds at the end of the slice if the current object is
	// a slice.
	//
	// Returns a Value pointing to that item.
	Append(value interface{}) Value

	// InsertAt inserts at the index-th position in the slice.
	//
	// Returns the slice itself, so that the calls can be chained,
	// or nil if the index is out-of-bound.
	InsertAt(index int, value interface{}) Slice
}

// Value is the most generic representation of the data.
type Value interface {
	// Data returns the underlying data provided by the object.
	Data() interface{}

	// Parent returns the parent interface or nil if no parent is
	// present (root).
	Parent() Value

	// Set changes the current value.
	//
	// Returns the Value itself, so that the calls can be chained.
	Set(value interface{}) Value

	// Map tries to convert the Value to a Map type. This will
	// return nil if it can't be converted.
	//
	// If the current value is nil, a new map will be created, set
	// and returned for convenience.
	Map() Map

	// Slice tries to convert the Value into a Slice type. This will
	// return nil if it can't be converted.
	//
	// If the current value is nil, a new slice will be created, set
	// and returned for conveninence.
	Slice() Slice
}

// New creates a representation of the given Value.
func New(value interface{}) Value {
	return &baseValue{data: value}
}

// Root returns the root for a given Value.
func Root(v Value) Value {
	for v.Parent() != nil {
		v = v.Parent()
	}
	return v
}
