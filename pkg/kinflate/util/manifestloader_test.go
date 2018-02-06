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

package util_test

import (
	"reflect"
	"testing"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

func TestManifestLoader(t *testing.T) {
	manifest := &manifest.Manifest{
		NamePrefix: "prefix",
	}
	loader := kutil.ManifestLoader{FS: fs.MakeFakeFS()}

	if err := loader.Write("my-manifest.yaml", manifest); err != nil {
		t.Fatalf("Couldn't write manifest file: %v\n", err)
	}

	readManifest, err := loader.Read("my-manifest.yaml")
	if err != nil {
		t.Fatalf("Couldn't read manifest file: %v\n", err)
	}
	if !reflect.DeepEqual(manifest, readManifest) {
		t.Fatal("Read manifest is different from written manifest")
	}
}

func TestManifestLoaderEmptyFile(t *testing.T) {
	manifest := &manifest.Manifest{
		NamePrefix: "prefix",
	}
	loader := kutil.ManifestLoader{}
	if loader.Write("", manifest) == nil {
		t.Fatalf("Write to empty filename should fail")
	}
}
