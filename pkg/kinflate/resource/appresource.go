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

package resource

import (
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/loader"
)

func resourcesFromPath(loader loader.Loader, path string) ([]*Resource, error) {
	content, err := loader.Load(path)
	if err != nil {
		return nil, err
	}

	objs, err := kutil.Decode(content)
	if err != nil {
		return nil, err
	}

	var res []*Resource
	for _, obj := range objs {
		res = append(res, &Resource{Data: obj})
	}
	return res, nil
}

//  NewFromPaths returns a slice of Resources given a resource path slice from manifest file.
func NewFromPaths(loader loader.Loader, paths []string) ([]*Resource, error) {
	allResources := []*Resource{}
	for _, path := range paths {
		res, err := resourcesFromPath(loader, path)
		if err != nil {
			return nil, err
		}
		allResources = append(allResources, res...)
	}
	return allResources, nil
}
