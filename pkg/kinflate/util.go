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
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"

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
func loadBaseAndOverlayPkg(f string) ([]string, []string, *manifest.Manifest, error) {
	overlay, err := loadManifestPkg(path.Join(f, kubeManifestFileName))
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: support `recursive` when we figure out what its behavior should be.
	// Recursive: overlay.Recursive
	overlayFiles := []string{}

	for _, o := range overlay.Patches {
		overlayFiles = append(overlayFiles, path.Join(f, o))
	}

	if len(overlay.Resources) == 0 {
		return nil, nil, nil, errors.New("expect at least one base, but got 0")
	}

	var baseFiles []string
	for _, base := range overlay.Resources {
		baseManifest, err := loadManifestPkg(path.Join(f, base, kubeManifestFileName))
		if err != nil {
			return nil, nil, nil, err
		}
		for _, filename := range baseManifest.Resources {
			baseFiles = append(baseFiles, path.Join(f, base, filename))
		}
	}

	return baseFiles, overlayFiles, overlay, nil
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
