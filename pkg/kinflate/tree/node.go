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

type namePrefixType string

type objectLabelsType map[string]string

type objectAnnotationsType map[string]string

type resourcesType types.KObject

type patchesType types.KObject

type configmapsType types.KObject

type secretsType types.KObject

// ManifestNode groups (possibly empty) manifest data with a (possibly empty)
// set of manifest nodes.
// data in one node may refer to data in other nodes.
// The node is invalid if it requires data it cannot find.
// A node is either a base app, or a patch to the base, or a patch to a patch to the base, etc.
type ManifestNode struct {
	data     *manifestData
	children []*ManifestNode
}

// manifestData contains all the objects loaded from the filesystem according to
// the Manifest Object.
type manifestData struct {
	name              string
	namePrefix        namePrefixType
	objectLabels      objectLabelsType
	objectAnnotations objectAnnotationsType
	resources         resourcesType
	patches           patchesType
	configmaps        configmapsType
	secrets           secretsType
}
