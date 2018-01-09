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
	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
)

// updateMetadata will inject the labels and annotations and add name prefix.
func updateMetadata(jsonObj []byte, overlayPkg *manifest.Manifest) ([]byte, error) {
	if len(jsonObj) == 0 || overlayPkg == nil {
		return nil, nil
	}

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(jsonObj, nil, nil)
	if err != nil {
		return nil, err
	}

	return updateObjectMetadata(obj, overlayPkg)
}

func updateObjectMetadata(obj runtime.Object, overlayPkg *manifest.Manifest) ([]byte, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	accessor.SetName(overlayPkg.NamePrefix + accessor.GetName())

	labels := accessor.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	for k, v := range overlayPkg.ObjectLabels {
		labels[k] = v
	}
	accessor.SetLabels(labels)

	annotations := accessor.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	for k, v := range overlayPkg.ObjectAnnotations {
		annotations[k] = v
	}
	accessor.SetAnnotations(annotations)

	return yaml.Marshal(obj)
}
