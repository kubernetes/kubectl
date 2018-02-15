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

package util

import (
	"fmt"
	"testing"
)

func TestManifestError_Error(t *testing.T) {
	filepath := "/path/to/Kube-manifest.yaml"
	errorMsg := "Manifest not found"
	expectedErrorMsg := fmt.Sprintf("Manifest File [%s]: %s\n", filepath, errorMsg)
	me := ManifestError{ManifestFilepath: filepath, ErrorMsg: errorMsg}
	if me.Error() != expectedErrorMsg {
		t.Errorf("Incorrect ManifestError.Error() message\n")
		t.Errorf("  Expected: %s\n", expectedErrorMsg)
		t.Errorf("  Got: %s\n", me.Error())
	}
}
