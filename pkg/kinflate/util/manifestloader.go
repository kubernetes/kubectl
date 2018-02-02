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
	"path"

	"github.com/ghodss/yaml"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

type ManifestLoader struct {
	FS fs.FileSystem
}

func (m *ManifestLoader) fs() fs.FileSystem {
	if m.FS == nil {
		m.FS = fs.MakeRealFS()
	}
	return m.FS
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
