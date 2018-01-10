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

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	kapps "k8s.io/kubectl/pkg/kinflate/apps"
	ptvisitor "k8s.io/kubectl/pkg/kinflate/pod_template_visitor"
)

func mutateField(m map[string]interface{}, pathToField []string, fn func(interface{}) (interface{}, error)) error {
	if len(pathToField) == 0 {
		return nil
	}

	v, found := m[pathToField[0]]
	if !found {
		return nil
	}

	if len(pathToField) == 1 {
		var err error
		m[pathToField[0]], err = fn(m[pathToField[0]])
		if err != nil {
			return err
		}
		return nil
	}

	newPathToField := pathToField[1:]
	switch typedV := v.(type) {
	case map[string]interface{}:
		return mutateField(typedV, newPathToField, fn)
	case []interface{}:
		for i := range typedV {
			item := typedV[i]
			typedItem, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%#v is expectd to be %T", item, typedItem)
			}
			err := mutateField(typedItem, newPathToField, fn)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%#v is not expected to be a primitive type", typedV)
	}
}

func changeNameAccordingToMapAndAddPrefix(m map[groupVersionKindName]newNameObject, gvk schema.GroupVersionKind, prefix string) func(interface{}) (interface{}, error) {
	return func(in interface{}) (interface{}, error) {
		s, ok := in.(string)
		if !ok {
			return nil, fmt.Errorf("%#v is expectd to be %T", in, s)
		}
		gvkn := groupVersionKindName{
			gvk:  gvk,
			name: s,
		}
		newNameObject, found := m[gvkn]
		if !found {
			return nil, fmt.Errorf("failed to find %#v in %#v", gvkn, m)
		}
		return prefix + newNameObject.newName, nil
	}
}

func addPrefix(prefix string) func(interface{}) (interface{}, error) {
	return func(in interface{}) (interface{}, error) {
		s, ok := in.(string)
		if !ok {
			return nil, fmt.Errorf("%#v is expectd to be %T", in, s)
		}
		return prefix + s, nil
	}
}

func addMap(additionalMap map[string]string) func(interface{}) (interface{}, error) {
	return func(in interface{}) (interface{}, error) {
		m, ok := in.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%#v is expectd to be %T", in, m)
		}
		for k, v := range additionalMap {
			m[k] = v
		}
		return m, nil
	}
}

func updatePodTemplateSpecMetadata(obj *unstructured.Unstructured, overlayPkg *manifest.Manifest, gvknToNewNameObject map[groupVersionKindName]newNameObject) error {
	pts := &ptvisitor.PodTemplateSpecVisitor{
		Object: obj,
		MungeFn: func(podTemplateSpec map[string]interface{}) error {
			err := mutateField(podTemplateSpec, []string{"labels"}, addMap(overlayPkg.Labels))
			if err != nil {
				return err
			}
			err = mutateField(podTemplateSpec, []string{"annotations"}, addMap(overlayPkg.Annotations))
			if err != nil {
				return err
			}

			err = updatePodSpecForConfigMap(podTemplateSpec, overlayPkg, gvknToNewNameObject)
			if err != nil {
				return err
			}
			return updatePodSpecForSecret(podTemplateSpec, overlayPkg, gvknToNewNameObject)
		},
	}

	unstructuredContent := obj.UnstructuredContent()
	gv := unstructuredContent["apiVersion"].(string)
	kind := unstructuredContent["kind"].(string)
	groupVersion, err := schema.ParseGroupVersion(gv)
	if err != nil {
		pts.Err = err
	}
	gkElement := kapps.GroupKindElement{
		GroupKind: schema.GroupKind{
			Group: groupVersion.Group,
			Kind:  kind,
		},
		IgnoreNonWorkloadError: true,
	}
	gkElement.Accept(pts)
	return pts.Err
}

// update the configmap name in spec.volumes.configMap.name and
// spec.containers.env.configMapKeyRef.
func updatePodSpecForConfigMap(podTemplateSpec map[string]interface{}, overlayPkg *manifest.Manifest, gvknToNewNameObject map[groupVersionKindName]newNameObject) error {
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Kind:    "ConfigMap",
	}
	return updatePodSpecBasedOnFieldKey(podTemplateSpec, "configMap", gvk, overlayPkg, gvknToNewNameObject)
}

// update the secret name in spec.volumes.secret.name and
// spec.containers.env.secretKeyRef.
func updatePodSpecForSecret(podTemplateSpec map[string]interface{}, overlayPkg *manifest.Manifest, gvknToNewNameObject map[groupVersionKindName]newNameObject) error {
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Kind:    "Secret",
	}
	return updatePodSpecBasedOnFieldKey(podTemplateSpec, "secret", gvk, overlayPkg, gvknToNewNameObject)
}

func updatePodSpecBasedOnFieldKey(
	podTemplateSpec map[string]interface{},
	fieldKey string,
	gvk schema.GroupVersionKind,
	overlayPkg *manifest.Manifest,
	gvknToNewNameObject map[groupVersionKindName]newNameObject) error {
	err := mutateField(podTemplateSpec, []string{"spec", "volumes", fieldKey, "name"}, changeNameAccordingToMapAndAddPrefix(gvknToNewNameObject, gvk, overlayPkg.NamePrefix))
	if err != nil {
		return err
	}
	err = mutateField(podTemplateSpec, []string{"spec", "containers", "env", "valueFrom", fieldKey + "KeyRef", "name"}, changeNameAccordingToMapAndAddPrefix(gvknToNewNameObject, gvk, overlayPkg.NamePrefix))
	if err != nil {
		return err
	}
	return mutateField(podTemplateSpec, []string{"spec", "containers", "envFrom", fieldKey + "Ref", "name"}, changeNameAccordingToMapAndAddPrefix(gvknToNewNameObject, gvk, overlayPkg.NamePrefix))
}

// updateMetadata will inject the labels and annotations and add name prefix.
func updateMetadata(jsonObj []byte, overlayPkg *manifest.Manifest, gvknToNewNameObject map[groupVersionKindName]newNameObject) ([]byte, error) {
	if len(jsonObj) == 0 || overlayPkg == nil {
		return nil, nil
	}

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(jsonObj, nil, nil)
	if err != nil {
		return nil, err
	}

	err = updatePodTemplateSpecMetadata(obj.(*unstructured.Unstructured), overlayPkg, gvknToNewNameObject)
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

	typeAccessor, err := meta.TypeAccessor(obj)
	if err != nil {
		return nil, err
	}
	if typeAccessor.GetAPIVersion() == "v1" && typeAccessor.GetKind() == "Service" {
		err = updateServiceSelector(obj.(*unstructured.Unstructured), overlayPkg.ObjectLabels)
		if err != nil {
			return nil, err
		}
	}

	return yaml.Marshal(obj)
}

func updateServiceSelector(obj *unstructured.Unstructured, labels map[string]string) error {
	objMap := obj.UnstructuredContent()
	return mutateField(objMap, []string{"spec", "selector"}, addMap(labels))
}
