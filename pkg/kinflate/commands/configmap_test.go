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

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

func TestNewAddConfigMapIsNotNil(t *testing.T) {
	if NewCmdAddConfigMap(nil, fs.MakeFakeFS()) == nil {
		t.Fatal("NewCmdAddConfigMap shouldn't be nil")
	}
}

func TestNewAddConfigMap_addConfigMap(t *testing.T) {

	testCases := []struct {
		testName              string
		config                dataConfig
		numExpectedConfigmaps int
		shouldFail            bool
	}{
		{
			testName: "first-time",
			config: dataConfig{
				Name: "test-config-name",
			},
			numExpectedConfigmaps: 1,
			shouldFail:            false,
		},
	}

	for _, test := range testCases {
		testManifest := manifest.Manifest{
			NamePrefix: "test-name-prefix",
		}
		// First time adding configmap to manifest
		err := addConfigMap(&testManifest, test.config)
		if err != nil && test.shouldFail {
			t.Fatal("Add configmap should not return error")
		}
		if test.numExpectedConfigmaps != len(testManifest.Configmaps) {
			t.Fatal("Manifest.Configmaps should have one entry after addConfigMap()")
		}
	}
}
