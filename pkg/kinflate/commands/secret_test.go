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

package commands

import (
	"testing"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

func TestNewAddSecretIsNotNil(t *testing.T) {
	if newCmdAddSecret(nil, fs.MakeFakeFS()) == nil {
		t.Fatal("newCmdAddSecret shouldn't be nil")
	}
}

func TestTLSecretExists(t *testing.T) {
	tlsName := "test-tls-secret"

	m := &manifest.Manifest{
		NamePrefix: "test-name-prefix",
	}

	if len(m.TLSSecrets) != 0 {
		t.Fatal("Initial manifest should not have any TLS secrets")
	}
	if tlsSecretExists(m, tlsName) {
		t.Fatalf("TLS Secret should not exist in manifest")
	}

	m.TLSSecrets = append(m.TLSSecrets, manifest.TLSSecret{Name: tlsName})

	if len(m.TLSSecrets) != 1 {
		t.Fatal("Manifest should have one TLS secrets")
	}
	if !tlsSecretExists(m, tlsName) {
		t.Fatalf("One TLS Secret should exist in manifest")
	}
}
