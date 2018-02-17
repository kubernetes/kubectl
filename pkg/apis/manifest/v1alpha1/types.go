/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Descriptor contains all the metadata of the package and drives package
// searching and browsing, and support the fork/rebase upgrade workflow.
// It can be used by something like an app registry.
type Descriptor struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// The name of the package.
	// The name of an individual instance should live in metadata.name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Description of the package.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// An pointer to the icon.
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`

	// Search keywords for the package.
	Keywords []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`

	// Homepage of the application package.
	Home string `json:"home,omitempty" yaml:"home,omitempty"`

	// Source specifies the upstream URL, e.g. https://github.com/foo/bar.git,
	// file://host/path, etc.
	// hosting the resource files specified in Base and Overlays.
	// This is useful in the fork/rebase workflow.
	Source string `json:"source,omitempty" yaml:"source,omitempty"`

	// Version of the package.
	PackageVersion string `json:"packageVersion,omitempty" yaml:"packageVersion,omitempty"`
}

// Manifest has all the information to expand of generate the k8s api resources.
// It can be used by kubectl or some other tooling.
// A manifest could be either a Base or an Overlay.
type Manifest struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// TODO: figure out if we need field ManifestVersion.
	// See: https://github.com/kubernetes/kubernetes/pull/52570/files/3eea91793dfbc3fdb0799589fac3790c4cde58a4#r140391019
	// Version of the manifest.
	// ManifestVersion string `json:"manifestVersion,omitempty" yaml:"manifestVersion,omitempty"`

	// NamePrefix will prefix the names of all resources mentioned in the manifest
	// including generated configmaps and secrets.
	NamePrefix string `json:"namePrefix,omitempty" yaml:"namePrefix,omitempty"`

	// Labels to add to all objects and selectors.
	// These labels would also be used to form the selector for apply --prune
	// Named differently than “labels” to avoid confusion with metadata for
	// this object
	ObjectLabels map[string]string `json:"objectLabels,omitempty" yaml:"objectLabels,omitempty"`

	// Annotations to add to all objects.
	ObjectAnnotations map[string]string `json:"objectAnnotations,omitempty" yaml:"objectAnnotations,omitempty"`

	// Packages contain the paths to other packages that this manifest depends on.
	// Each path should be either a path to a Kube-manifest.yaml or a path of
	// a directory that contains a Kube-manifest.yaml file.
	Packages []string `json:"packages,omitempty" yaml:"packages,omitempty"`

	// Resources specifies the relative paths within the package.
	// It could be any format that kubectl -f allows, i.e. files, directories,
	// URLs and globs.
	Resources []string `json:"resources,omitempty" yaml:"resources,omitempty"`

	// An Patch entry is very similar to an Resource entry.
	// It specifies the relative paths within the package, and could be any
	// format that kubectl -f allows.
	// It should be able to be merged by Strategic Merge Patch on top of its
	// corresponding base resource.
	Patches []string `json:"patches,omitempty" yaml:"patches,omitempty"`

	// List of configmaps to generate from configuration sources.
	// Base/overlay concept doesn't apply to this field.
	// If a configmap want to have a base and an overlay, it should go to Bases
	// and Overlays fields.
	Configmaps []ConfigMap `json:"configmaps,omitempty" yaml:"configmaps,omitempty"`

	// List of secrets to generate from secret commands.
	// Base/overlay concept doesn't apply to this field.
	// If a secret want to have a base and an overlay, it should go to Bases and
	// Overlays fields.
	SecretGenerators []SecretGenerator `json:"secretGenerators,omitempty" yaml:"secretGenerators,omitempty"`

	// Whether prune resources not defined in Kube-manifest.yaml, similar to
	// `kubectl apply --prune` behavior.
	Prune bool `json:"prune,omitempty" yaml:"prune,omitempty"`

	// TODO: figure out what the behavior details should be.
	// Whether PersistentVolumeClaims should be deleted with the other resources.
	// OwnPersistentVolumeClaims bool `json:"ownPersistentVolumeClaims,omitempty" yaml:"ownPersistentVolumeClaims,omitempty"`

	// TODO: figure out what the behavior details should be.
	// Whether recursive look for Kube-manifest.yaml, similar to
	// `kubectl --recursive` behavior.
	// Recursive bool `json:"recursive,omitempty" yaml:"recursive,omitempty"`
}

// ConfigMap contains the metadata of how to generate a configmap.
type ConfigMap struct {
	// Name of the configmap.
	// The full name should be Manifest.NamePrefix + Configmap.Name +
	// hash(content of configmap).
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// DataSources for configmap.
	DataSources `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// SecretGenerator contains the metadata of how to generate a secret.
type SecretGenerator struct {
	// Name of the secret.
	// The full name should be Manifest.NamePrefix + SecretGenerator.Name +
	// hash(content of secret).
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Type of the secret.
	//
	// This is the same field as the secret type field in v1/Secret:
	// It can be "Opaque" (default), or "kubernetes.io/tls".
	//
	// If type is "kubernetes.io/tls", then "Commands" must have exactly two
	// keys: "tls.key" and "tls.crt"
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Map of keys to commands to generate the values
	Commands map[string]string `json:",commands,omitempty" yaml:",inline,omitempty"`
}

// DataSources contains some generic sources for configmap or secret.
// Only one field can be set.
type DataSources struct {
	// LiteralSources is a list of literal sources.
	// Each literal source should be a key and literal value,
	// e.g. `somekey=somevalue`
	// It will be similar to kubectl create configmap|secret --from-literal
	LiteralSources []string `json:"literals,omitempty" yaml:"literals,omitempty"`

	// FileSources is a list of file sources.
	// Each file source can be specified using its file path, in which case file
	// basename will be used as configmap key, or optionally with a key and file
	// path, in which case the given key will be used.
	// Specifying a directory will iterate each named file in the directory
	// whose basename is a valid configmap key.
	// It will be similar to kubectl create configmap|secret --from-file
	FileSources []string `json:"files,omitempty" yaml:"files,omitempty"`

	// EnvSource format should be a path to a file to read lines of key=val
	// pairs to create a configmap.
	// i.e. a Docker .env file or a .ini file.
	EnvSource string `json:"env,omitempty" yaml:"env,omitempty"`
}
