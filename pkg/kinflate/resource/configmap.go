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

package resource

import (
	corev1 "k8s.io/api/core/v1"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	cutil "k8s.io/kubectl/pkg/kinflate/configmapandsecret/util"
)

// (Note): pass in loader which has rootPath context and knows how to load
// files given a relative path.
// NewFromConfigMap returns a Resource given a configmap metadata from manifest
// file.
func NewFromConfigMap(cm manifest.ConfigMap) (*Resource, error) {
	corev1CM, err := makeConfigMap(cm)
	if err != nil {
		return nil, err
	}

	data, err := objectToUnstructured(corev1CM)
	if err != nil {
		return nil, err
	}
	return &Resource{Data: data}, nil
}

func makeConfigMap(cm manifest.ConfigMap) (*corev1.ConfigMap, error) {
	corev1cm := &corev1.ConfigMap{}
	corev1cm.APIVersion = "v1"
	corev1cm.Kind = "ConfigMap"
	corev1cm.Name = cm.Name
	corev1cm.Data = map[string]string{}

	// TODO: move the configmap helpers functions in this file/package
	if cm.EnvSource != "" {
		if err := cutil.HandleConfigMapFromEnvFileSource(corev1cm, cm.EnvSource); err != nil {
			return nil, err
		}
	}
	if cm.FileSources != nil {
		if err := cutil.HandleConfigMapFromFileSources(corev1cm, cm.FileSources); err != nil {
			return nil, err
		}
	}
	if cm.LiteralSources != nil {
		if err := cutil.HandleConfigMapFromLiteralSources(corev1cm, cm.LiteralSources); err != nil {
			return nil, err
		}
	}

	return corev1cm, nil
}
