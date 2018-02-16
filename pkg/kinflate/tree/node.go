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
	"fmt"

	"k8s.io/kubectl/pkg/kinflate/transformers"
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

func (md *ManifestData) allResources() error {
	err := types.Merge(md.Resources, md.Configmaps)
	if err != nil {
		return err
	}
	return types.Merge(md.Resources, md.Secrets)
}

// ModeType is the option type for kinflate inflate
type ModeType string

const (
	// ModeNormal means regular transformation.
	ModeNormal ModeType = "normal_mode"
	// ModeNoop means no transformation.
	ModeNoop ModeType = "noop_mode"
)

func (md *ManifestData) preprocess(mode ModeType) error {
	switch mode {
	case ModeNormal:
		return md.allResources()
	case ModeNoop:
		return nil
	default:
		return fmt.Errorf("unknown mode for inflate")
	}
}

func (md *ManifestData) makeTransformer(mode ModeType) (transformers.Transformer, error) {
	switch mode {
	case ModeNormal:
		return DefaultTransformer(md)
	case ModeNoop:
		return transformers.NewNoOpTransformer(), nil
	default:
		return transformers.NewNoOpTransformer(), fmt.Errorf("unknown mode for inflate")
	}
}

// Inflate will recursively do the transformation on all the nodes below.
func (md *ManifestData) Inflate(mode ModeType) error {
	for _, pkg := range md.Packages {
		err := pkg.Inflate(ModeNormal)
		if err != nil {
			return err
		}
	}

	for _, pkg := range md.Packages {
		err := types.Merge(md.Resources, pkg.Resources)
		if err != nil {
			return err
		}
	}

	err := md.preprocess(mode)
	if err != nil {
		return err
	}

	t, err := md.makeTransformer(mode)
	if err != nil {
		return err
	}
	return t.Transform(types.KObject(md.Resources))
}
