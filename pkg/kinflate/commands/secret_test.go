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

package commands

import (
	"testing"
)

func TestNewAddSecretIsNotNil(t *testing.T) {
	if NewCmdAddSecret(nil) == nil {
		t.Fatal("NewCmdAddSecret shouldn't be nil")
	}
}

func TestAddGenericSecretValidation_NoName(t *testing.T) {
	config := addGenericSecret{}

	if config.Validate([]string{}) == nil {
		t.Fatal("Validation should fail if no name is specified")
	}
}

func TestAddGenericSecretValidation_MoreThanOneName(t *testing.T) {
	config := addGenericSecret{}

	if config.Validate([]string{"name", "othername"}) == nil {
		t.Fatal("Validation should fail if more than one name is specified")
	}
}

func TestAddGenericSecretValidation_Flags(t *testing.T) {
	tests := []struct {
		name       string
		config     addGenericSecret
		shouldFail bool
	}{
		{
			name: "env-file-source and literal are both set",
			config: addGenericSecret{
				LiteralSources: []string{"one", "two"},
				EnvFileSource:  "three",
			},
			shouldFail: true,
		},
		{
			name: "env-file-source and from-file are both set",
			config: addGenericSecret{
				FileSources:   []string{"one", "two"},
				EnvFileSource: "three",
			},
			shouldFail: true,
		},
		{
			name:       "we don't have any option set",
			config:     addGenericSecret{},
			shouldFail: true,
		},
		{
			name: "we have from-file and literal ",
			config: addGenericSecret{
				LiteralSources: []string{"one", "two"},
				FileSources:    []string{"three", "four"},
			},
			shouldFail: false,
		},
	}

	for _, test := range tests {
		if test.config.Validate([]string{"name"}) == nil && test.shouldFail {
			t.Fatalf("Validation should fail if %s", test.name)
		} else if test.config.Validate([]string{"name"}) != nil && !test.shouldFail {
			t.Fatalf("Validation should succeed if %s", test.name)
		}
	}
}

func TestAddTLSSecretValidation_NoName(t *testing.T) {
	config := addTLSSecret{}

	if config.Validate([]string{}) == nil {
		t.Fatal("Validation should fail if no name is specified")
	}
}

func TestAddTLSSecretValidation_MoreThanOneName(t *testing.T) {
	config := addTLSSecret{}

	if config.Validate([]string{"name", "othername"}) == nil {
		t.Fatal("Validation should fail if more than one name is specified")
	}
}

func TestAddTLSSecretValidation_Flags(t *testing.T) {
	tests := []struct {
		name       string
		config     addTLSSecret
		shouldFail bool
	}{
		{
			name: "cert and key are set",
			config: addTLSSecret{
				Cert: "cert",
				Key:  "key",
			},
			shouldFail: false,
		},
		{
			name: "cert is set, but not key",
			config: addTLSSecret{
				Cert: "cert",
			},
			shouldFail: true,
		},
		{
			name: "key is set, but not cert",
			config: addTLSSecret{
				Key: "key",
			},
			shouldFail: true,
		},
		{
			name:       "neither key nor cert is set",
			config:     addTLSSecret{},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		if test.config.Validate([]string{"name"}) == nil && test.shouldFail {
			t.Fatalf("Validation should fail if %s", test.name)
		} else if test.config.Validate([]string{"name"}) != nil && !test.shouldFail {
			t.Fatalf("Validation should succeed if %s", test.name)
		}
	}
}
