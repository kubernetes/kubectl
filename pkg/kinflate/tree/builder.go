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
	"os"
	"path/filepath"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	cutil "k8s.io/kubectl/pkg/kinflate/configmapandsecret"
	"k8s.io/kubectl/pkg/kinflate/types"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

type Loader struct {
	FS   fs.FileSystem
	Path string
}

// LoadManifestDataFromPath takes a path to a Kube-manifest.yaml or a dir that has a Kube-manifest.yaml.
// It returns a tree of ManifestData.
func (l *Loader) LoadManifestDataFromPath(path string) (*ManifestData, error) {
	return l.loadManifestDataFromPath(path)
}

// loadManifestDataFromPath make a ManifestData from path
func (l *Loader) loadManifestDataFromPath(path string) (*ManifestData, error) {
	m, err := l.loadManifestFileFromPath(path)
	if err != nil {
		return nil, err
	}
	return l.manifestToManifestData(m)
}

// manifestToManifestData make a ManifestData given an Manifest object
func (l *Loader) manifestToManifestData(m *manifest.Manifest) (*ManifestData, error) {
	mdata, err := l.loadManifestDataFromManifestFileAndResources(m)
	if err != nil {
		return nil, err
	}

	pkgs := []*ManifestData{}
	for _, pkg := range m.Packages {
		pkgNode, err := l.LoadManifestDataFromPath(pkg)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, pkgNode)
	}
	mdata.Packages = pkgs
	return mdata, nil
}

// loadManifestFileFromPath loads a manifest object from file.
func (l *Loader) loadManifestFileFromPath(filename string) (*manifest.Manifest, error) {
	var loader = kutil.ManifestLoader{FS: l.FS}
	fullFilename, err := loader.MakeValidManifestPath(filename)
	if err != nil {
		return nil, err
	}
	l.Path = fullFilename
	m, err := loader.Read(fullFilename)
	if err != nil {
		return nil, err
	}
	loader.Validate(m)
	return m, err
}

func (l *Loader) loadManifestDataFromManifestFileAndResources(m *manifest.Manifest) (*ManifestData, error) {
	mdata := &ManifestData{}
	var err error
	mdata.Name = m.Name
	mdata.NamePrefix = NamePrefixType(m.NamePrefix)
	mdata.ObjectLabels = m.ObjectLabels
	mdata.ObjectAnnotations = m.ObjectAnnotations

	res, err := l.LoadKObjectFromPaths(m.Resources)
	if err != nil {
		errorMsg := fmt.Sprintf("Resource from Manifest (%s) couldn't be loaded properly.\n%v\n"+
			"Please check the Resource subsection in (%s).", l.Path, err, l.Path)
		return nil, fmt.Errorf(errorMsg)
	}
	mdata.Resources = ResourcesType(res)

	pat, err := l.LoadKObjectFromPaths(m.Patches)
	if err != nil {
		return nil, err
	}
	mdata.Patches = PatchesType(pat)

	cms, err := cutil.MakeConfigMapsKObject(m.Configmaps)
	if err != nil {
		return nil, err
	}
	mdata.Configmaps = ConfigmapsType(cms)

	sec, err := cutil.MakeGenericSecretsKObject(m.GenericSecrets)
	if err != nil {
		return nil, err
	}
	mdata.Secrets = SecretsType(sec)

	TLS, err := cutil.MakeTLSSecretsKObject(m.TLSSecrets)
	err = types.Merge(mdata.Secrets, TLS)
	if err != nil {
		return nil, err
	}
	return mdata, nil
}

func (l *Loader) LoadKObjectFromPaths(paths []string) (types.KObject, error) {
	res := types.KObject{}
	for _, path := range paths {
		err := l.LoadKObjectFromPath(path, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (l *Loader) LoadKObjectFromPath(path string, into types.KObject) error {
	_, err := l.FS.Stat(path)
	if err != nil {
		return err
	}
	if into == nil {
		return fmt.Errorf("cannot load object to an empty KObject")
	}

	var e error
	filepath.Walk(path, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			e = err
			return err
		}
		// Skip all the dir
		if info.IsDir() {
			return nil
		}

		err = l.LoadKObjectFromFile(filepath, into)
		return nil
	})
	return e
}

func (l *Loader) LoadKObjectFromFile(filename string, into types.KObject) error {
	f, err := l.FS.Stat(filename)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return fmt.Errorf("%q is NOT expected to be an dir", filename)
	}
	content, err := l.FS.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = kutil.Decode(content, into)
	if err != nil {
		return err
	}
	return nil
}
