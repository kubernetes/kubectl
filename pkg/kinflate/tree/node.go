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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/kinflate/gvkn"
)

type NamePrefixType string

type ObjectLabelsType map[string]string

type ObjectAnnotationsType map[string]string

type ResourcesType map[gvkn.GroupVersionKindName]*unstructured.Unstructured

type PatchesType map[gvkn.GroupVersionKindName]*unstructured.Unstructured

type ConfigmapsType map[gvkn.GroupVersionKindName]*unstructured.Unstructured

type SecretsType map[gvkn.GroupVersionKindName]*unstructured.Unstructured

// ManifestNode is the node for building the manifest tree.
// Children points to the packages defined on this node's Manifest.
type ManifestNode struct {
	Data     *ManifestData
	Children []*ManifestNode
}

// ManifestData contains all the objects loaded from the filesystem according to
// the Manifest Object.
type ManifestData struct {
	Name              string
	NamePrefix        NamePrefixType
	ObjectLabels      ObjectLabelsType
	ObjectAnnotations ObjectAnnotationsType
	Resources         ResourcesType
	Patches           PatchesType
	Configmaps        ConfigmapsType
	Secrets           SecretsType
}
