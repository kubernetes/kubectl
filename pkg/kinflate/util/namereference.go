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

package util

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NameReferenceTransformer contains the referencing info between 2 GroupVersionKinds
type NameReferenceTransformer struct {
	pathConfigs []referencePathConfig
}

var _ Transformer = &NameReferenceTransformer{}

// NewDefaultingNameReferenceTransformer constructs a NameReferenceTransformer
// with defaultNameReferencepathConfigs.
func NewDefaultingNameReferenceTransformer() (*NameReferenceTransformer, error) {
	return NewNameReferenceTransformer(defaultNameReferencePathConfigs)
}

// NewNameReferenceTransformer construct a NameReferenceTransformer.
func NewNameReferenceTransformer(pc []referencePathConfig) (*NameReferenceTransformer, error) {
	if pc == nil {
		return nil, errors.New("pathConfigs is not expected to be nil")
	}
	return &NameReferenceTransformer{pathConfigs: pc}, nil
}

// Transform does the fields update according to pathConfigs.
// The old name is in the key in the map and the new name is in the object
// associated with the key. e.g. if <k, v> is one of the key-value pair in the map,
// then the old name is k.Name and the new name is v.GetName()
func (o *NameReferenceTransformer) Transform(
	m map[GroupVersionKindName]*unstructured.Unstructured) error {
	for GVKn := range m {
		obj := m[GVKn]
		objMap := obj.UnstructuredContent()
		for _, referencePathConfig := range o.pathConfigs {
			for _, path := range referencePathConfig.pathConfigs {
				if !SelectByGVK(GVKn.GVK, path.GroupVersionKind) {
					continue
				}
				err := mutateField(objMap, path.Path, path.CreateIfNotPresent,
					o.updateNameReference(referencePathConfig.referencedGVK, m))
				// Ignore the error when we can't find the GVKN that is being
				// referenced, because the missing GVKN may be not included in
				// this manifest and will be created later.
				if IsNoMatchingGVKNError(err) {
					continue
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// NoMatchingGVKNError indicates failing to find a GroupVersionKindName.
type NoMatchingGVKNError struct {
	message string
}

// NewNoMatchingGVKNError constructs an instance of NoMatchingGVKNError with
// a given error message.
func NewNoMatchingGVKNError(errMsg string) NoMatchingGVKNError {
	return NoMatchingGVKNError{errMsg}
}

// IsNoMatchingGVKNError checks if the error is NoMatchingGVKNError type.
func IsNoMatchingGVKNError(err error) bool {
	_, ok := err.(NoMatchingGVKNError)
	return ok
}

// Error returns the error in string format.
func (err NoMatchingGVKNError) Error() string {
	return err.message
}

func (o *NameReferenceTransformer) updateNameReference(
	GVK schema.GroupVersionKind,
	m map[GroupVersionKindName]*unstructured.Unstructured,
) func(in interface{}) (interface{}, error) {
	return func(in interface{}) (interface{}, error) {
		s, ok := in.(string)
		if !ok {
			return nil, fmt.Errorf("%#v is expectd to be %T", in, s)
		}

		for GVKn, obj := range m {
			if !SelectByGVK(GVKn.GVK, &GVK) {
				continue
			}
			if GVKn.Name == s {
				return obj.GetName(), nil
			}
		}
		return nil, NewNoMatchingGVKNError(
			fmt.Sprintf("no matching for GroupVersionKind %v and Name %v", GVK, s))
	}
}
