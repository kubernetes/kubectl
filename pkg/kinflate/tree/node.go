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

package tree

import (
	"k8s.io/kubectl/pkg/kinflate/types"
)

type NamePrefixType string

type ObjectLabelsType map[string]string

type ObjectAnnotationsType map[string]string

type ResourcesType types.KObject

type PatchesType types.KObject

type ConfigmapsType types.KObject

type SecretsType types.KObject

type PackagesType []*ManifestData

// ManifestData contains all the objects loaded from the filesystem according to
// the Manifest Object.
// Data in one node may refer to data in other nodes.
// The node is invalid if it requires data it cannot find.
// A node is either a base app, or a patch to the base, or a patch to a patch to the base, etc.
type ManifestData struct {
	// Name of the manifest
	Name string

	NamePrefix        NamePrefixType
	ObjectLabels      ObjectLabelsType
	ObjectAnnotations ObjectAnnotationsType
	Resources         ResourcesType
	Patches           PatchesType
	Configmaps        ConfigmapsType
	Secrets           SecretsType

	Packages PackagesType
}

func (md *ManifestData) AllResources() error {
	err := types.Merge(md.Resources, md.Configmaps)
	if err != nil {
		return err
	}
	return types.Merge(md.Resources, md.Secrets)
}

func (md *ManifestData) Inflate() error {
	for _, pkg := range md.Packages {
		err := pkg.Inflate()
		if err != nil {
			return err
		}
		err = types.Merge(md.Resources, pkg.Resources)
		if err != nil {
			return err
		}
	}

	err := md.AllResources()
	if err != nil {
		return err
	}
	t, err := DefaultTransformer(md)
	return t.Transform(types.KObject(md.Resources))
}
