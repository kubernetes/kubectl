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

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/scheme"
)

const kubeManifestFileName = "Kube-manifest.yaml"

// loadManifestPkg loads a manifest file and parse it in to the Package object.
func loadManifestPkg(dirname string) (*manifest.Manifest, error) {
	bytes, err := ioutil.ReadFile(path.Join(dirname, kubeManifestFileName))
	if err != nil {
		return nil, err
	}
	var manifest manifest.Manifest
	// TODO: support json
	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}
	adjustPathsForManifest(&manifest, []string{dirname})
	return &manifest, err
}

func populateResourceMap(files []string,
	m map[kutil.GroupVersionKindName]*unstructured.Unstructured) error {
	for _, file := range files {
		_, err := fileToMap(file, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileToMap(filename string,
	into map[kutil.GroupVersionKindName]*unstructured.Unstructured,
) (map[kutil.GroupVersionKindName]*unstructured.Unstructured, error) {
	f, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if into == nil {
		into = map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		_, err = DirToMap(filename, into)
	case mode.IsRegular():
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		_, err = kutil.Decode(content, into)
		if err != nil {
			return nil, err
		}
	}
	return into, nil
}

// DirToMap tries to find Kube-manifest.yaml first in a dir.
// If not found, traverse all the file in the dir.
func DirToMap(dirname string,
	into map[kutil.GroupVersionKindName]*unstructured.Unstructured,
) (map[kutil.GroupVersionKindName]*unstructured.Unstructured, error) {
	if into == nil {
		into = map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	}
	f, err := os.Stat(dirname)
	if !f.IsDir() {
		return nil, fmt.Errorf("%q is expected to be an dir", dirname)
	}

	kubeManifestFileAbsName := path.Join(dirname, kubeManifestFileName)
	_, err = os.Stat(kubeManifestFileAbsName)
	switch {
	case err != nil && !os.IsNotExist(err):
		return nil, err
	case err == nil:
		manifest, err := loadManifestPkg(dirname)
		if err != nil {
			return nil, err
		}
		_, err = ManifestToMap(manifest, into)
		if err != nil {
			return nil, err
		}
	case err != nil && os.IsNotExist(err):
		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			_, err = fileToMap(path.Join(dirname, file.Name()), into)
		}

		var e error
		filepath.Walk(dirname, func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				e = err
				return err
			}
			_, err = fileToMap(path, into)
			return nil
		})
	}
	return into, nil
}

// ManifestToMap builds a map given the info in manifest.
func ManifestToMap(manifest *manifest.Manifest,
	into map[kutil.GroupVersionKindName]*unstructured.Unstructured,
) (map[kutil.GroupVersionKindName]*unstructured.Unstructured, error) {
	baseResouceMap := map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	if into != nil {
		baseResouceMap = into
	}
	err := populateResourceMap(manifest.Resources, baseResouceMap)
	if err != nil {
		return nil, err
	}

	overlayResouceMap := map[kutil.GroupVersionKindName]*unstructured.Unstructured{}
	err = populateResourceMap(manifest.Patches, overlayResouceMap)
	if err != nil {
		return nil, err
	}

	// Strategic merge the resources exist in both base and overlay.
	for gvkn, base := range baseResouceMap {
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
			baseResouceMap[gvkn].Object = merged
			delete(overlayResouceMap, gvkn)
		}
	}

	// If there are resources in overlay that are not defined in base, just add it to base.
	if len(overlayResouceMap) > 0 {
		for gvkn, jsonObj := range overlayResouceMap {
			baseResouceMap[gvkn] = jsonObj
		}
	}

	err = populateConfigMapAndSecretMap(manifest, baseResouceMap)
	if err != nil {
		return nil, err
	}

	transformers, err := DefaultTransformers(manifest)
	if err != nil {
		return nil, err
	}
	for _, transformer := range transformers {
		err = transformer.Transform(baseResouceMap)
		if err != nil {
			return nil, err
		}
	}
	return baseResouceMap, nil
}

// DefaultTransformers generates 4 transformers:
// 1) name prefix 2) apply labels 3) apply annotations 4) update name reference
func DefaultTransformers(manifest *manifest.Manifest) ([]kutil.Transformer, error) {
	transformers := []kutil.Transformer{}

	npt, err := kutil.NewDefaultingNamePrefixTransformer(manifest.NamePrefix)
	if err != nil {
		return nil, err
	}
	if npt != nil {
		transformers = append(transformers, npt)
	}

	lt, err := kutil.NewDefaultingLabelsMapTransformer(manifest.ObjectLabels)
	if err != nil {
		return nil, err
	}
	if lt != nil {
		transformers = append(transformers, lt)
	}

	at, err := kutil.NewDefaultingAnnotationsMapTransformer(manifest.ObjectAnnotations)
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
