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

package transformers

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/kinflate/gvkn"
)

// MultiTransformer contains a list of transformers.
type MultiTransformer struct {
	Transformers []Transformer
}

var _ Transformer = &MultiTransformer{}

// Transform prepends the name prefix.
func (o *MultiTransformer) Transform(m map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	for _, t := range o.Transformers {
		err := t.Transform(m)
		if err != nil {
			return err
		}
	}
	return nil
}
