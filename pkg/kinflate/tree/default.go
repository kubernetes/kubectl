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

package tree

import (
	"k8s.io/kubectl/pkg/kinflate/transformers"
	"k8s.io/kubectl/pkg/kinflate/types"
)

// DefaultTransformer generates the following transformers:
// 1) apply overlay
// 2) name prefix
// 3) apply labels
// 4) apply annotations
// 5) update name reference
func DefaultTransformer(m *ManifestData) (transformers.Transformer, error) {
	ts := []transformers.Transformer{}

	ot, err := transformers.NewOverlayTransformer(types.KObject(m.Patches))
	if err != nil {
		return nil, err
	}
	if ot != nil {
		ts = append(ts, ot)
	}

	npt, err := transformers.NewDefaultingNamePrefixTransformer(string(m.NamePrefix))
	if err != nil {
		return nil, err
	}
	if npt != nil {
		ts = append(ts, npt)
	}

	lt, err := transformers.NewDefaultingLabelsMapTransformer(m.ObjectLabels)
	if err != nil {
		return nil, err
	}
	if lt != nil {
		ts = append(ts, lt)
	}

	at, err := transformers.NewDefaultingAnnotationsMapTransformer(m.ObjectAnnotations)
	if err != nil {
		return nil, err
	}
	if at != nil {
		ts = append(ts, at)
	}

	nrt, err := transformers.NewDefaultingNameReferenceTransformer()
	if err != nil {
		return nil, err
	}
	if nrt != nil {
		ts = append(ts, nrt)
	}
	return transformers.NewMultiTransformer(ts), nil
}
