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

package mergemap

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/kinflate/types"
)

// Merge will merge all the entries in m2 to m1.
func Merge(m1, m2 map[types.GroupVersionKindName]*unstructured.Unstructured,
) error {
	for gvkn, obj := range m2 {
		if _, found := m1[gvkn]; found {
			return fmt.Errorf("there is already an entry: %q", gvkn)
		}
		m1[gvkn] = obj
	}
	return nil
}
