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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/scheme"
)

func populateResourceMap(files []string,
	m map[kutil.GroupVersionKindName]*unstructured.Unstructured) error {
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
) (map[kutil.GroupVersionKindName]*unstructured.Unstructured, error) {
	f, err := os.Stat(mPath)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		mPath = path.Join(mPath, KubeManifestFileName)
	} else {
		if !strings.HasSuffix(mPath, KubeManifestFileName) {
			return nil, fmt.Errorf("expecting file: %q, but got: %q", KubeManifestFileName, mPath)
		}
	}
	manifest, err := (&kutil.ManifestLoader{}).Read(mPath)
	if err != nil {
		return nil, err
	}
	return ManifestToMap(manifest)
}

func pathToMap(path string, into map[kutil.GroupVersionKindName]*unstructured.Unstructured) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if into == nil {
		into = map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		err = dirToMap(path, into)
	case mode.IsRegular():
		err = fileToMap(path, into)
	}
	return err
}

func fileToMap(filename string, into map[kutil.GroupVersionKindName]*unstructured.Unstructured) error {
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
func dirToMap(dirname string, into map[kutil.GroupVersionKindName]*unstructured.Unstructured) error {
	if into == nil {
		into = map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	}
	f, err := os.Stat(dirname)
	if !f.IsDir() {
		return fmt.Errorf("%q is expected to be an dir", dirname)
	}

	kubeManifestFileAbsName := path.Join(dirname, KubeManifestFileName)
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
) (map[kutil.GroupVersionKindName]*unstructured.Unstructured, error) {
	return manifestToMap(m, nil)
}

// manifestToMap takes a manifest and recursively finds all instances of Kube-manifest,
// reads them and merges them all into `into`.
func manifestToMap(m *manifest.Manifest,
	into map[kutil.GroupVersionKindName]*unstructured.Unstructured,
) (map[kutil.GroupVersionKindName]*unstructured.Unstructured, error) {
	baseResourceMap := map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	if into != nil {
		baseResourceMap = into
	}
	err := populateResourceMap(m.Resources, baseResourceMap)
	if err != nil {
		return nil, err
	}

	overlayResouceMap := map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
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
			merged, err := strategicpatch.StrategicMergeMapPatch(
				base.UnstructuredContent(),
				overlay.UnstructuredContent(),
				versionedObj)
			if err != nil {
				return nil, err
			}
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

	transformers, err := defaultTransformers(m)
	if err != nil {
		return nil, err
	}
	for _, t := range transformers {
		err = t.Transform(baseResourceMap)
		if err != nil {
			return nil, err
		}
	}
	return baseResourceMap, nil
}

// defaultTransformers generates 4 transformers:
// 1) name prefix 2) apply labels 3) apply annotations 4) update name reference
func defaultTransformers(m *manifest.Manifest) ([]kutil.Transformer, error) {
	transformers := []kutil.Transformer{}

	npt, err := kutil.NewDefaultingNamePrefixTransformer(m.NamePrefix)
	if err != nil {
		return nil, err
	}
	if npt != nil {
		transformers = append(transformers, npt)
	}

	lt, err := kutil.NewDefaultingLabelsMapTransformer(m.ObjectLabels)
	if err != nil {
		return nil, err
	}
	if lt != nil {
		transformers = append(transformers, lt)
	}

	at, err := kutil.NewDefaultingAnnotationsMapTransformer(m.ObjectAnnotations)
	if err != nil {
		return nil, err
	}
	if at != nil {
		transformers = append(transformers, at)
	}

	nrt, err := kutil.NewDefaultingNameReferenceTransformer()
	if err != nil {
		return nil, err
	}
	if nrt != nil {
		transformers = append(transformers, nrt)
	}
	return transformers, nil
}
