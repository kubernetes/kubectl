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

type PrefixNameOptions struct {
	prefix      string
	pathsConfig *PathsConfig
}

var _ Transformer = &PrefixNameOptions{}

var DefaultNamePrefixPathsConfig = &PathsConfig{
	[]PathConfig{
		{
			Path:               []string{"metadata", "name"},
			CreateIfNotPresent: false,
		},
	},
}

func (o *PrefixNameOptions) Complete(prefix string, pathsConfig *PathsConfig) {
	o.prefix = prefix
	if pathsConfig == nil {
		pathsConfig = DefaultNamePrefixPathsConfig
	}
	o.pathsConfig = pathsConfig
}

func (o *PrefixNameOptions) Transform(m map[GroupVersionKindName]*unstructured.Unstructured) error {
	for gvkn := range m {
		obj := m[gvkn]
		objMap := obj.UnstructuredContent()
		for _, path := range o.pathsConfig.Paths {
			err := mutateField(objMap, path.Path, path.CreateIfNotPresent, o.addPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *PrefixNameOptions) TransformBytes(in []byte) ([]byte, error) {
	m, err := Decode(in)
	if err != nil {
		return nil, err
	}

	err = o.Transform(m)
	if err != nil {
		return nil, err
	}

	return Encode(m)
}

func (o *PrefixNameOptions) addPrefix(in interface{}) (interface{}, error) {
	s, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("%#v is expectd to be %T", in, s)
	}
	return o.prefix + s, nil
}
