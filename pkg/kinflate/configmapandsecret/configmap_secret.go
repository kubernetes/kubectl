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

package configmapandsecret

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	cutil "k8s.io/kubectl/pkg/kinflate/configmapandsecret/util"
	"k8s.io/kubectl/pkg/kinflate/hash"
	"k8s.io/kubectl/pkg/kinflate/types"
)

// MakeConfigmapAndGenerateName makes a configmap and returns the configmap and the name appended with a hash.
func MakeConfigmapAndGenerateName(cm manifest.ConfigMap) (*unstructured.Unstructured, string, error) {
	corev1CM, err := makeConfigMap(cm)
	if err != nil {
		return nil, "", err
	}
	h, err := hash.ConfigMapHash(corev1CM)
	if err != nil {
		return nil, "", err
	}
	nameWithHash := fmt.Sprintf("%s-%s", corev1CM.GetName(), h)
	unstructuredCM, err := objectToUnstructured(corev1CM)
	return unstructuredCM, nameWithHash, err
}

// MakeSecretAndGenerateName returns a secret with the name appended with a hash.
func MakeSecretAndGenerateName(secret manifest.SecretGenerator, path string) (*unstructured.Unstructured, string, error) {
	corev1Secret, err := makeSecret(secret, path)
	if err != nil {
		return nil, "", err
	}
	h, err := hash.SecretHash(corev1Secret)
	if err != nil {
		return nil, "", err
	}
	nameWithHash := fmt.Sprintf("%s-%s", secret.Name, h)
	unstructuredCM, err := objectToUnstructured(corev1Secret)
	return unstructuredCM, nameWithHash, err
}

func objectToUnstructured(in runtime.Object) (*unstructured.Unstructured, error) {
	marshaled, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var out unstructured.Unstructured
	err = out.UnmarshalJSON(marshaled)
	return &out, err
}

func makeConfigMap(cm manifest.ConfigMap) (*corev1.ConfigMap, error) {
	corev1cm := &corev1.ConfigMap{}
	corev1cm.APIVersion = "v1"
	corev1cm.Kind = "ConfigMap"
	corev1cm.Name = cm.Name
	corev1cm.Data = map[string]string{}

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

func makeSecret(secret manifest.SecretGenerator, path string) (*corev1.Secret, error) {
	corev1secret := &corev1.Secret{}
	corev1secret.APIVersion = "v1"
	corev1secret.Kind = "Secret"
	corev1secret.Name = secret.Name
	corev1secret.Type = corev1.SecretType(secret.Type)
	if corev1secret.Type == "" {
		corev1secret.Type = corev1.SecretTypeOpaque
	}
	corev1secret.Data = map[string][]byte{}

	for k, v := range secret.Commands {
		out, err := createSecretKey(path, v)
		if err != nil {
			return nil, err
		}
		corev1secret.Data[k] = out
	}

	return corev1secret, nil
}

func populateMap(m types.KObject, obj *unstructured.Unstructured, newName string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	oldName := accessor.GetName()
	gvk := obj.GetObjectKind().GroupVersionKind()
	gvkn := types.GroupVersionKindName{GVK: gvk, Name: oldName}

	if _, found := m[gvkn]; found {
		return fmt.Errorf("The <name: %q, GroupVersionKind: %v> already exists in the map", oldName, gvk)
	}
	accessor.SetName(newName)
	m[gvkn] = obj
	return nil
}

// MakeConfigMapsKObject returns a map of <GVK, oldName> -> unstructured object.
func MakeConfigMapsKObject(maps []manifest.ConfigMap) (types.KObject, error) {
	m := types.KObject{}
	for _, cm := range maps {
		unstructuredConfigMap, nameWithHash, err := MakeConfigmapAndGenerateName(cm)
		if err != nil {
			return nil, err
		}
		err = populateMap(m, unstructuredConfigMap, nameWithHash)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// MakeSecretsKObject returns a map of <GVK, oldName> -> unstructured object.
func MakeSecretsKObject(secrets []manifest.SecretGenerator, path string) (types.KObject, error) {
	m := types.KObject{}
	for _, secret := range secrets {
		unstructuredSecret, nameWithHash, err := MakeSecretAndGenerateName(secret, path)
		if err != nil {
			return nil, err
		}
		err = populateMap(m, unstructuredSecret, nameWithHash)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func createSecretKey(wd string, command string) ([]byte, error) {
	fi, err := os.Stat(wd)
	if err != nil || !fi.IsDir() {
		wd = filepath.Dir(wd)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = wd

	return cmd.Output()
}
