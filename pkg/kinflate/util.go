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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	cutil "k8s.io/kubectl/pkg/kinflate/configmapandsecret"
	"k8s.io/kubectl/pkg/kinflate/constants"
	"k8s.io/kubectl/pkg/kinflate/gvkn"
	"k8s.io/kubectl/pkg/kinflate/mergemap"
	"k8s.io/kubectl/pkg/kinflate/transformers"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/scheme"
)

func populateMap(m map[gvkn.GroupVersionKindName]*unstructured.Unstructured, obj *unstructured.Unstructured, newName string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	oldName := accessor.GetName()
	gvk := obj.GetObjectKind().GroupVersionKind()
	gvkn := gvkn.GroupVersionKindName{GVK: gvk, Name: oldName}

	if _, found := m[gvkn]; found {
		return fmt.Errorf("cannot use a duplicate name %q for %s", oldName, gvk)
	}
	accessor.SetName(newName)
	m[gvkn] = obj
	return nil
}

func populateConfigMapAndSecretMap(manifest *manifest.Manifest, m map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	configmaps, err := cutil.MakeMapOfConfigMap(manifest)
	if err != nil {
		return err
	}
	err = mergemap.Merge(m, configmaps)
	if err != nil {
		return err
	}

	genericSecrets, err := cutil.MakeMapOfGenericSecret(manifest)
	if err != nil {
		return err
	}
	err = mergemap.Merge(m, genericSecrets)
	if err != nil {
		return err
	}

	TLSSecrets, err := cutil.MakeMapOfTLSSecret(manifest)
	if err != nil {
		return err
	}
	return mergemap.Merge(m, TLSSecrets)
}

func populateResourceMap(files []string,
	m map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	for _, file := range files {
		err := pathToMap(file, m)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadFromManifestPath loads the manifest from the given path.
// It returns a map of resources defined in the manifest file.
func LoadFromManifestPath(mPath string,
) (map[gvkn.GroupVersionKindName]*unstructured.Unstructured, error) {
	f, err := os.Stat(mPath)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		mPath = path.Join(mPath, constants.KubeManifestFileName)
	} else {
		if !strings.HasSuffix(mPath, constants.KubeManifestFileName) {
			return nil, fmt.Errorf("expecting file: %q, but got: %q", constants.KubeManifestFileName, mPath)
		}
	}
	manifest, err := (&kutil.ManifestLoader{}).Read(mPath)
	if err != nil {
		return nil, err
	}
	return ManifestToMap(manifest)
}

func pathToMap(path string, into map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if into == nil {
		into = map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		err = dirToMap(path, into)
	case mode.IsRegular():
		err = fileToMap(path, into)
	}
	return err
}

func fileToMap(filename string, into map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	f, err := os.Stat(filename)
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

// dirToMap tries to find Kube-manifest.yaml first in a dir.
// If not found, traverse all the file in the dir.
func dirToMap(dirname string, into map[gvkn.GroupVersionKindName]*unstructured.Unstructured) error {
	if into == nil {
		into = map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	}
	f, err := os.Stat(dirname)
	if !f.IsDir() {
		return fmt.Errorf("%q is expected to be an dir", dirname)
	}

	kubeManifestFileAbsName := path.Join(dirname, constants.KubeManifestFileName)
	_, err = os.Stat(kubeManifestFileAbsName)
	switch {
	case err != nil && !os.IsNotExist(err):
		return err
	case err == nil:
		manifest, err := (&kutil.ManifestLoader{}).Read(kubeManifestFileAbsName)
		if err != nil {
			return err
		}
		_, err = manifestToMap(manifest, into)
		if err != nil {
			return err
		}
	case err != nil && os.IsNotExist(err):
		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			return err
		}

		for _, file := range files {
			err = pathToMap(path.Join(dirname, file.Name()), into)
		}

		var e error
		filepath.Walk(dirname, func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				e = err
				return err
			}
			err = fileToMap(path, into)
			return nil
		})
	}
	return nil
}

// ManifestToMap takes a manifest and recursively finds all instances of Kube-manifest,
// reads them and merges them all in a map of resources.
func ManifestToMap(m *manifest.Manifest,
) (map[gvkn.GroupVersionKindName]*unstructured.Unstructured, error) {
	return manifestToMap(m, nil)
}

// manifestToMap takes a manifest and recursively finds all instances of Kube-manifest,
// reads them and merges them all into `into`.
func manifestToMap(m *manifest.Manifest,
	into map[gvkn.GroupVersionKindName]*unstructured.Unstructured,
) (map[gvkn.GroupVersionKindName]*unstructured.Unstructured, error) {
	baseResourceMap := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	if into != nil {
		baseResourceMap = into
	}
	err := populateResourceMap(m.Resources, baseResourceMap)
	if err != nil {
		return nil, err
	}

	overlayResouceMap := map[gvkn.GroupVersionKindName]*unstructured.Unstructured{}
	err = populateResourceMap(m.Patches, overlayResouceMap)
	if err != nil {
		return nil, err
	}

	// Strategic merge the resources exist in both base and overlay.
	for gvkn, base := range baseResourceMap {
		// Merge overlay with base resource.
		if overlay, found := overlayResouceMap[gvkn]; found {
			versionedObj, err := scheme.Scheme.New(gvkn.GVK)
			if err != nil {
				switch {
				case runtime.IsNotRegisteredError(err):
					return nil, fmt.Errorf("CRD and TPR are not supported now: %v", err)
				default:
					return nil, err
				}
			}
			// Store the name of the base object, because this name may have been munged.
			// Apply this name to the StrategicMergePatched object.
			baseName := base.GetName()
			merged, err := strategicpatch.StrategicMergeMapPatch(
				base.UnstructuredContent(),
				overlay.UnstructuredContent(),
				versionedObj)
			if err != nil {
				return nil, err
			}
			base.SetName(baseName)
			baseResourceMap[gvkn].Object = merged
			delete(overlayResouceMap, gvkn)
		}
	}

	// If there are resources in overlay that are not defined in base, just add it to base.
	if len(overlayResouceMap) > 0 {
		for gvkn, jsonObj := range overlayResouceMap {
			baseResourceMap[gvkn] = jsonObj
		}
	}

	err = populateConfigMapAndSecretMap(m, baseResourceMap)
	if err != nil {
		return nil, err
	}

	t, err := transformers.DefaultTransformer(m)
	if err != nil {
		return nil, err
	}
	err = t.Transform(baseResourceMap)
	if err != nil {
		return nil, err
	}
	return baseResourceMap, nil
}
