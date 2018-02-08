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
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/adjustpath"
	cutil "k8s.io/kubectl/pkg/kinflate/configmapandsecret"
	"k8s.io/kubectl/pkg/kinflate/constants"
	"k8s.io/kubectl/pkg/kinflate/gvkn"
	"k8s.io/kubectl/pkg/kinflate/mergemap"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
)

// BuildManifestTree takes a path to a Kube-manifest.yaml or a dir that has a Kube-manifest.yaml.
// It returns a tree of ManifestNode.
func BuildManifestTree(path string) (*ManifestNode, error) {
	return manifestPathToManifestNode(path)
}

func manifestPathToManifestNode(path string) (*ManifestNode, error) {
	path, err := validateManifestPath(path)
	if err != nil {
		return nil, err
	}
	m, err := manifestPathToManifest(path)
	if err != nil {
		return nil, err
	}
	return manifestToManifestNode(m)
}

func manifestToManifestNode(m *manifest.Manifest) (*ManifestNode, error) {
	mnode := &ManifestNode{}
	var err error
	mnode.Data, err = manifestToManifestData(m)
	if err != nil {
		return nil, err
	}

	mnode.Children = []*ManifestNode{}
	for _, pkg := range m.Packages {
		child, err := manifestPathToManifestNode(pkg)
		if err != nil {
			return nil, err
		}
		mnode.Children = append(mnode.Children, child)
	}
	return mnode, nil
}

// manifestPathToManifest loads a manifest file and parse it in to the Manifest object.
func manifestPathToManifest(filename string) (*manifest.Manifest, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var m manifest.Manifest
	// TODO: support json
	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}
	dir, _ := path.Split(filename)
	adjustpath.AdjustPathsForManifest(&m, []string{dir})
	return &m, err
}

// validateManifestPath loads the manifest from the given path.
// It returns ManifestData and an potential error.
func validateManifestPath(mPath string) (string, error) {
	f, err := os.Stat(mPath)
	if err != nil {
		return "", err
	}
	if f.IsDir() {
		mPath = path.Join(mPath, constants.KubeManifestFileName)
		_, err = os.Stat(mPath)
		if err != nil {
			return "", err
		}
	} else {
		if !strings.HasSuffix(mPath, constants.KubeManifestFileName) {
			return "", fmt.Errorf("expecting file: %q, but got: %q", constants.KubeManifestFileName, mPath)
		}
	}
	return mPath, nil
}

func manifestToManifestData(m *manifest.Manifest) (*ManifestData, error) {
	mdata := &ManifestData{}
	var err error
	mdata.Name = m.Name
	mdata.NamePrefix = NamePrefixType(m.NamePrefix)
	mdata.ObjectLabels = m.ObjectLabels
	mdata.ObjectAnnotations = m.ObjectAnnotations
	mdata.Resources, err = pathsToMap(m.Resources)
	if err != nil {
		return nil, err
	}
	mdata.Patches, err = pathsToMap(m.Patches)
	if err != nil {
		return nil, err
	}
	mdata.Configmaps, err = cutil.MakeMapOfConfigMap(m)
	if err != nil {
		return nil, err
	}
	mdata.Secrets, err = cutil.MakeMapOfGenericSecret(m)
	if err != nil {
		return nil, err
	}
	TLSSecrets, err := cutil.MakeMapOfTLSSecret(m)
	err = mergemap.Merge(mdata.Secrets, TLSSecrets)
	if err != nil {
		return nil, err
	}
	return mdata, nil
}

func pathsToMap(paths []string) (map[gvkn.GroupVersionKindName]*unstructured.Unstructured, error) {
	res := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	for _, path := range paths {
		err := pathToMap(path, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func pathToMap(path string, into map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	if into == nil {
		into = map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
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

		err = fileToMap(filepath, into)
		return nil
	})
	return e
}

func fileToMap(filename string, into map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
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
