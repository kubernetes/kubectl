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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"

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

// MakeGenericSecretAndGenerateName makes a generic secret and returns the secret and the name appended with a hash.
func MakeGenericSecretAndGenerateName(secret manifest.GenericSecret) (*unstructured.Unstructured, string, error) {
	corev1Secret, err := makeGenericSecret(secret)
	if err != nil {
		return nil, "", err
	}
	return makeSecretAndGenerateName(corev1Secret, secret.Name)
}

// MakeTLSSecretAndGenerateName makes a generic secret and returns the secret and the name appended with a hash.
func MakeTLSSecretAndGenerateName(secret manifest.TLSSecret) (*unstructured.Unstructured, string, error) {
	corev1Secret, err := makeTlsSecret(secret)
	if err != nil {
		return nil, "", err
	}
	return makeSecretAndGenerateName(corev1Secret, secret.Name)
}

func makeSecretAndGenerateName(secret *corev1.Secret, name string) (*unstructured.Unstructured, string, error) {
	h, err := hash.SecretHash(secret)
	if err != nil {
		return nil, "", err
	}
	nameWithHash := fmt.Sprintf("%s-%s", name, h)
	unstructuredCM, err := objectToUnstructured(secret)
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

func makeGenericSecret(secret manifest.GenericSecret) (*corev1.Secret, error) {
	corev1secret := &corev1.Secret{}
	corev1secret.APIVersion = "v1"
	corev1secret.Kind = "Secret"
	corev1secret.Name = secret.Name
	corev1secret.Type = corev1.SecretTypeOpaque
	corev1secret.Data = map[string][]byte{}

	if secret.EnvSource != "" {
		if err := cutil.HandleFromEnvFileSource(corev1secret, secret.EnvSource); err != nil {
			return nil, err
		}
	}
	if secret.FileSources != nil {
		if err := cutil.HandleFromFileSources(corev1secret, secret.FileSources); err != nil {
			return nil, err
		}
	}
	if secret.LiteralSources != nil {
		if err := cutil.HandleFromLiteralSources(corev1secret, secret.LiteralSources); err != nil {
			return nil, err
		}
	}
	return corev1secret, nil
}

func makeTlsSecret(secret manifest.TLSSecret) (*corev1.Secret, error) {
	corev1secret := &corev1.Secret{}
	corev1secret.APIVersion = "v1"
	corev1secret.Kind = "Secret"
	corev1secret.Name = secret.Name
	corev1secret.Type = corev1.SecretTypeTLS
	corev1secret.Data = map[string][]byte{}

	if err := validateTLS(secret.CertFile, secret.KeyFile); err != nil {
		return nil, err
	}
	tlsCrt, err := ioutil.ReadFile(secret.CertFile)
	if err != nil {
		return nil, err
	}
	tlsKey, err := ioutil.ReadFile(secret.KeyFile)
	if err != nil {
		return nil, err
	}
	corev1secret.Data[corev1.TLSCertKey] = []byte(tlsCrt)
	corev1secret.Data[corev1.TLSPrivateKeyKey] = []byte(tlsKey)

	return corev1secret, err
}

func validateTLS(cert, key string) error {
	if len(key) == 0 {
		return fmt.Errorf("key must be specified")
	}
	if len(cert) == 0 {
		return fmt.Errorf("certificate must be specified")
	}
	if _, err := tls.LoadX509KeyPair(cert, key); err != nil {
		return fmt.Errorf("failed to load key pair %v", err)
	}
	return nil
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

// MakeGenericSecretsKObject returns a map of <GVK, oldName> -> unstructured object.
func MakeGenericSecretsKObject(secrets []manifest.GenericSecret) (types.KObject, error) {
	m := types.KObject{}
	for _, secret := range secrets {
		unstructuredSecret, nameWithHash, err := MakeGenericSecretAndGenerateName(secret)
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

// MakeTLSSecretsKObject returns a map of <GVK, oldName> -> unstructured object.
func MakeTLSSecretsKObject(secrets []manifest.TLSSecret) (types.KObject, error) {
	m := types.KObject{}
	for _, secret := range secrets {
		unstructuredSecret, nameWithHash, err := MakeTLSSecretAndGenerateName(secret)
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
