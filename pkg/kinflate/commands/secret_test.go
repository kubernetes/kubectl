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
)

func TestNewAddSecretIsNotNil(t *testing.T) {
	if NewCmdAddSecret(nil) == nil {
		t.Fatal("NewCmdAddSecret shouldn't be nil")
	}
}

func TestGetOrCreateGenericSecret(t *testing.T) {
	gsName := "test-generic-secret"

	manifest := &manifest.Manifest{
		NamePrefix: "test-name-prefix",
	}

	if len(manifest.GenericSecrets) != 0 {
		t.Fatal("Initial manifest should not have any genericsecrets")
	}
	gs := getOrCreateGenericSecret(manifest, gsName)

	if gs == nil {
		t.Fatalf("GenericSecret should always be non-nil")
	}

	if len(manifest.GenericSecrets) != 1 {
		t.Fatalf("Manifest should have newly created generic secret")
	}

	if &manifest.GenericSecrets[len(manifest.GenericSecrets)-1] != gs {
		t.Fatalf("Pointer address for newly inserted generic secret should be same")
	}

	existingGS := getOrCreateGenericSecret(manifest, gsName)

	if existingGS != gs {
		t.Fatalf("should have returned an existing generic secret with name: %v", gsName)
	}

	if len(manifest.GenericSecrets) != 1 {
		t.Fatalf("Should not insert generic secret for an existing name: %v", gsName)
	}
}
