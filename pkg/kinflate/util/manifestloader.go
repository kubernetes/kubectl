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

package util

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/ghodss/yaml"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/constants"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

type ManifestLoader struct {
	FS fs.FileSystem
}

// First pass to encapsulate fields for more informative error messages.
type ManifestErrors struct {
	filepath string
	errorMsg string
}

func (m *ManifestLoader) fs() fs.FileSystem {
	if m.FS == nil {
		m.FS = fs.MakeRealFS()
	}
	return m.FS
}

// makeValidManifestPath returns a path to a KubeManifest file known to exist.
// The argument is either the full path to the file itself, or a path to a directory
// that immediately contains the file. Anything else is an error.
func (m *ManifestLoader) MakeValidManifestPath(mPath string) (string, error) {
	f, err := m.fs().Stat(mPath)
	if err != nil {
		errorMsg := fmt.Sprintf("Manifest (%s) missing\nRun `kinflate init` first", mPath)
		return "", errors.New(errorMsg)
	}
	if f.IsDir() {
		mPath = path.Join(mPath, constants.KubeManifestFileName)
		_, err = m.fs().Stat(mPath)
		if err != nil {
			errorMsg := fmt.Sprintf("Manifest (%s) missing\nRun `kinflate init` first", mPath)
			return "", errors.New(errorMsg)
		}
	} else {
		if !strings.HasSuffix(mPath, constants.KubeManifestFileName) {
			return "", fmt.Errorf("Manifest file (%s) should have %s suffix\n", mPath, constants.KubeManifestSuffix)
		}
	}
	return mPath, nil
}

// Read loads a manifest file and parse it in to the Manifest object.
func (m *ManifestLoader) Read(filename string) (*manifest.Manifest, error) {
	bytes, err := m.fs().ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var manifest manifest.Manifest
	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}
	dir, _ := path.Split(filename)
	adjustPathsForManifest(&manifest, []string{dir})
	return &manifest, err
}

// Write dumps the Manifest object into a file. If manifest is nil, an
// error is returned.
func (m *ManifestLoader) Write(filename string, manifest *manifest.Manifest) error {
	if manifest == nil {
		return errors.New("util: failed to write passed-in nil manifest")
	}
	bytes, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	return m.fs().WriteFile(filename, bytes)
}

// Read must have already been called and we have a loaded manifest
func (m *ManifestLoader) Validate(manifest *manifest.Manifest) []ManifestErrors {
	//TODO: implement this function
	//// validate Packages
	//merrors := m.validatePackages(manifest.Packages)
	//// validate Resources
	//merrors = merrors + m.validateResources(manifest.Resources)
	//
	//// validate Patches
	//merrors = append(merrors, m.validatePatches(manifest.Patches))
	//
	//// validate Configmaps
	//merrors = append(merrors, m.validateConfigmaps(manifest.Configmaps))
	//
	//// validate GenericSecrets
	//merrors = append(merrors, m.validateGenericSecrets(manifest.GenericSecrets))
	//
	//// validate TLSSecrets
	//merrors = append(merrors, m.validateTLSSecrets(manifest.TLSSecrets))
	return nil
}
