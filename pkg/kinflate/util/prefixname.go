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
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// PrefixNameOptions contains the prefix and the path config for each field that
// the name prefix will be applied.
type PrefixNameOptions struct {
	prefix      string
	pathConfigs []PathConfig
}

var _ Transformer = &PrefixNameOptions{}

var DefaultNamePrefixPathConfigs = []PathConfig{
	{
		Path:               []string{"metadata", "name"},
		CreateIfNotPresent: false,
	},
}

func (o *PrefixNameOptions) Complete(prefix string, pathConfigs []PathConfig) {
	o.prefix = prefix
	if pathConfigs == nil {
		pathConfigs = DefaultNamePrefixPathConfigs
	}
	o.pathConfigs = pathConfigs
}

func (o *PrefixNameOptions) Transform(m map[GroupVersionKindName]*unstructured.Unstructured) error {
	for gvkn := range m {
		obj := m[gvkn]
		objMap := obj.UnstructuredContent()
		for _, path := range o.pathConfigs {
			if !SelectByGVK(gvkn.GVK, path.GroupVersionKind) {
				continue
			}
			err := mutateField(objMap, path.Path, path.CreateIfNotPresent, o.addPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *PrefixNameOptions) addPrefix(in interface{}) (interface{}, error) {
	s, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("%#v is expectd to be %T", in, s)
	}
	return o.prefix + s, nil
}
