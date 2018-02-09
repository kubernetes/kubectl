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
	"io/ioutil"
	"os"
	"path/filepath"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	cutil "k8s.io/kubectl/pkg/kinflate/configmapandsecret"
	"k8s.io/kubectl/pkg/kinflate/mergemap"
	"k8s.io/kubectl/pkg/kinflate/types"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

// LoadManifestNodeFromPath takes a path to a Kube-manifest.yaml or a dir that has a Kube-manifest.yaml.
// It returns a tree of ManifestNode.
func LoadManifestNodeFromPath(path string) (*ManifestNode, error) {
	return loadManifestNodeFromPath(path)
}

// loadManifestNodeFromPath make a ManifestNode from path
func loadManifestNodeFromPath(path string) (*ManifestNode, error) {
	m, err := loadManifestFileFromPath(path)
	if err != nil {
		return nil, err
	}
	return manifestToManifestNode(m)
}

// manifestToManifestNode make a ManifestNode given an Manifest object
func manifestToManifestNode(m *manifest.Manifest) (*ManifestNode, error) {
	mnode := &ManifestNode{}
	var err error
	mnode.data, err = loadManifestDataFromManifestFileAndResources(m)
	if err != nil {
		return nil, err
	}

	mnode.children = []*ManifestNode{}
	for _, pkg := range m.Packages {
		child, err := loadManifestNodeFromPath(pkg)
		if err != nil {
			return nil, err
		}
		mnode.children = append(mnode.children, child)
	}
	return mnode, nil
}

// loadManifestFileFromPath loads a manifest object from file.
func loadManifestFileFromPath(filename string) (*manifest.Manifest, error) {
	loader := kutil.ManifestLoader{fs.MakeRealFS()}
	m, err := loader.Read(filename)
	if err != nil {
		return nil, err
	}
	return m, err
}

func loadManifestDataFromManifestFileAndResources(m *manifest.Manifest) (*manifestData, error) {
	mdata := &manifestData{}
	var err error
	mdata.name = m.Name
	mdata.namePrefix = namePrefixType(m.NamePrefix)
	mdata.objectLabels = m.ObjectLabels
	mdata.objectAnnotations = m.ObjectAnnotations

	res, err := loadKObjectFromPaths(m.Resources)
	if err != nil {
		return nil, err
	}
	mdata.resources = resourcesType(res)

	pat, err := loadKObjectFromPaths(m.Patches)
	if err != nil {
		return nil, err
	}
	mdata.patches = patchesType(pat)

	cms, err := cutil.MakeConfigMapsKObject(m.Configmaps)
	if err != nil {
		return nil, err
	}
	mdata.configmaps = configmapsType(cms)

	sec, err := cutil.MakeGenericSecretsKObject(m.GenericSecrets)
	if err != nil {
		return nil, err
	}
	mdata.secrets = secretsType(sec)

	TLS, err := cutil.MakeTLSSecretsKObject(m.TLSSecrets)
	err = mergemap.Merge(mdata.secrets, TLS)
	if err != nil {
		return nil, err
	}
	return mdata, nil
}

func loadKObjectFromPaths(paths []string) (types.KObject, error) {
	res := types.KObject{}
	for _, path := range paths {
		err := loadKObjectFromPath(path, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func loadKObjectFromPath(path string, into types.KObject) error {
	_, err := os.Stat(path)
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

		err = loadKObjectFromFile(filepath, into)
		return nil
	})
	return e
}

func loadKObjectFromFile(filename string, into types.KObject) error {
	f, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return fmt.Errorf("%q is NOT expected to be an dir", filename)
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = kutil.Decode(content, into)
	if err != nil {
		return err
	}
	return nil
}
