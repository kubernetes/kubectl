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
	"bytes"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubectl/pkg/kinflate/types"
)

// decode decodes a list of objects in byte array format
func decode(in []byte) ([]*unstructured.Unstructured, error) {
	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(in), 1024)
	objs := []*unstructured.Unstructured{}

	var err error
	for {
		var out unstructured.Unstructured
		err = decoder.Decode(&out)
		if err != nil {
			break
		}
		objs = append(objs, &out)
	}
	if err != io.EOF {
		return nil, err
	}
	return objs, nil
}

// decodeToResourceCollection decodes a list of objects in byte array format.
// it will return a ResourceCollection.
func decodeToResourceCollection(in []byte) (ResourceCollection, error) {
	objs, err := decode(in)
	if err != nil {
		return nil, err
	}

	into := ResourceCollection{}
	for i, obj := range objs {
		gvkn := types.GroupVersionKindName{
			GVK:  obj.GroupVersionKind(),
			Name: obj.GetName(),
		}
		if _, found := into[gvkn]; found {
			return into, fmt.Errorf("GroupVersionKindName: %#v already exists in the map", gvkn)
		}
		into[gvkn] = &Resource{Data: objs[i]}
	}
	return into, nil
}

func resourceCollectionFromResources(resources []*Resource) (ResourceCollection, error) {
	out := ResourceCollection{}
	for _, res := range resources {
		gvkn := res.GVKN()
		if _, found := out[gvkn]; found {
			return nil, fmt.Errorf("duplicated %#v is not allowed", gvkn)
		}
		out[gvkn] = res
	}
	return out, nil
}

// Merge will merge all of the entries in the slice of ResourceCollection.
func Merge(rcs ...ResourceCollection) (ResourceCollection, error) {
	all := ResourceCollection{}
	for _, rc := range rcs {
		for gvkn, obj := range rc {
			if _, found := all[gvkn]; found {
				return nil, fmt.Errorf("there is already an entry: %q", gvkn)
			}
			all[gvkn] = obj
		}
	}

	return all, nil
}
