/*
Copyright 2017 The Kubernetes Authors.

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

package transformers

import (
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
)

// DefaultTransformer generates 4 transformers:
// 1) name prefix 2) apply labels 3) apply annotations 4) update name reference
func DefaultTransformer(m *manifest.Manifest) (Transformer, error) {
	transformers := []Transformer{}

	npt, err := NewDefaultingNamePrefixTransformer(m.NamePrefix)
	if err != nil {
		return nil, err
	}
	if npt != nil {
		transformers = append(transformers, npt)
	}

	lt, err := NewDefaultingLabelsMapTransformer(m.ObjectLabels)
	if err != nil {
		return nil, err
	}
	if lt != nil {
		transformers = append(transformers, lt)
	}

	at, err := NewDefaultingAnnotationsMapTransformer(m.ObjectAnnotations)
	if err != nil {
		return nil, err
	}
	if at != nil {
		transformers = append(transformers, at)
	}

	nrt, err := NewDefaultingNameReferenceTransformer()
	if err != nil {
		return nil, err
	}
	if nrt != nil {
		transformers = append(transformers, nrt)
	}
	return &MultiTransformer{transformers}, nil
}
