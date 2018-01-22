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

package kinflate

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
)

const kubeManifestFileName = "Kube-manifest.yaml"

type resource struct {
	resources  []string
	configmaps []manifest.ConfigMap
	secrets    []manifest.Secret
}

type newNameObject struct {
	newName string
	obj     runtime.Object
}

// loadBaseAndOverlayPkg returns:
// - List of FilenameOptions, each FilenameOptions contains all the files and whether recursive for each base defined in overlay kube-manifest.yaml.
// - Fileoptions for overlay.
// - Package object for overlay.
// - A potential error.
func loadBaseAndOverlayPkg(f string) ([]*resource, *resource, *manifest.Manifest, error) {
	overlay, err := loadManifestPkg(path.Join(f, kubeManifestFileName))
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: support `recursive` when we figure out what its behavior should be.
	// Recursive: overlay.Recursive

	overlayResource := adjustPathsForConfigMapAndSecret(overlay, []string{f})
	patchResources, err := adjustPaths(overlay.Patches, []string{f})
	if err != nil {
		return nil, nil, nil, err
	}

	resources, err := adjustPaths(overlay.Resources, []string{f})
	if err != nil {
		return nil, nil, nil, err
	}
	overlayResource.resources = append(patchResources, resources...)

	if len(overlay.Resources) == 0 && len(overlay.Packages) == 0 {
		return nil, nil, nil, errors.New("expect at least one resource or one package, but got 0")
	}

	if err != nil {
		return nil, nil, nil, err
	}
	var baseResources []*resource

	for _, base := range overlay.Packages {
		baseManifest, err := loadManifestPkg(path.Join(f, base, kubeManifestFileName))
		if err != nil {
			return nil, nil, nil, err
		}

		baseResource := adjustPathsForConfigMapAndSecret(baseManifest, []string{f, base})
		baseResource.resources, err = adjustPaths(baseManifest.Resources, []string{f, base})
		if err != nil {
			return nil, nil, nil, err
		}

		baseResources = append(baseResources, baseResource)
	}

	return baseResources, overlayResource, overlay, nil
}

// loadManifestPkg loads a manifest file and parse it in to the Package object.
func loadManifestPkg(filename string) (*manifest.Manifest, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var pkg manifest.Manifest
	// TODO: support json
	err = yaml.Unmarshal(bytes, &pkg)
	return &pkg, err
}

func populateResourceMap(files []string, m map[groupVersionKindName][]byte, errOut io.Writer) error {
	decoder := unstructured.UnstructuredJSONScheme

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		// try converting to json, if there is a error, probably because the content is already json.
		jsoncontent, err := yaml.YAMLToJSON(content)
		if err != nil {
			fmt.Fprintf(errOut, "error when trying to convert yaml to json: %v\n", err)
		} else {
			content = jsoncontent
		}

		obj, gvk, err := decoder.Decode(content, nil, nil)
		if err != nil {
			return err
		}
		accessor, err := meta.Accessor(obj)
		if err != nil {
			return err
		}
		name := accessor.GetName()
		gvkn := groupVersionKindName{gvk: *gvk, name: name}
		if err != nil {
			return err
		}
		if _, found := m[gvkn]; found {
			return fmt.Errorf("unexpected same groupVersionKindName: %#v", gvkn)
		}
		m[gvkn] = content
	}
	return nil
}
