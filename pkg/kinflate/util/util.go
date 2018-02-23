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
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubectl/pkg/kinflate/types"
)

// Decode decodes a list of objects in byte array format
func Decode(in []byte) ([]*unstructured.Unstructured, error) {
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

// DecodeToKObject decodes a list of objects in byte array format.
// Decoded object will be inserted in `into` if it's not nil. Otherwise, it will
// construct a new map and return it.
func DecodeToKObject(in []byte, into types.KObject) (types.KObject, error) {
	objs, err := Decode(in)
	if err != nil {
		return nil, err
	}

	if into == nil {
		into = types.KObject{}
	}
	for i := range objs {
		metaAccessor, err := meta.Accessor(objs[i])
		if err != nil {
			return nil, err
		}
		name := metaAccessor.GetName()
		typeAccessor, err := meta.TypeAccessor(objs[i])
		if err != nil {
			return nil, err
		}
		apiVersion := typeAccessor.GetAPIVersion()
		kind := typeAccessor.GetKind()
		gv, err := schema.ParseGroupVersion(apiVersion)
		if err != nil {
			return nil, err
		}
		gvk := gv.WithKind(kind)
		gvkn := types.GroupVersionKindName{
			GVK:  gvk,
			Name: name,
		}
		if _, found := into[gvkn]; found {
			return into, fmt.Errorf("GroupVersionKindName: %#v already exists in the map", gvkn)
		}
		into[gvkn] = objs[i]
	}
	return into, nil
}

// Encode encodes the map `in` and output the encoded objects separated by `---`.
func Encode(in types.KObject) ([]byte, error) {
	gvknList := []types.GroupVersionKindName{}
	for gvkn := range in {
		gvknList = append(gvknList, gvkn)
	}
	sort.Sort(types.ByGVKN(gvknList))

	firstObj := true
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, gvkn := range gvknList {
		obj := in[gvkn]
		out, err := yaml.Marshal(obj)
		if err != nil {
			return nil, err
		}
		if !firstObj {
			_, err = buf.WriteString("---\n")
			if err != nil {
				return nil, err
			}
		}
		_, err = buf.Write(out)
		if err != nil {
			return nil, err
		}
		firstObj = false
	}
	return buf.Bytes(), nil
}

// WriteToDir write each object in KObject to a file named with GroupVersionKindName.
func WriteToDir(in types.KObject, dirName string, printer Printer) (*Directory, error) {
	dir, err := CreateDirectory(dirName)
	if err != nil {
		return &Directory{}, err
	}

	for gvkn, obj := range in {
		f, err := dir.NewFile(gvkn.String())
		if err != nil {
			return &Directory{}, err
		}
		defer f.Close()
		err = printer.Print(obj, f)
		if err != nil {
			return &Directory{}, err
		}
	}
	return dir, nil
}
