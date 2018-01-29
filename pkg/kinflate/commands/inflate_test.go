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
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const (
	input    = "../examples/simple/instances/exampleinstance/"
	expected = "testdata/simple/out/expected.yaml"
)

func TestInflate(t *testing.T) {
	const updateEnvVar = "UPDATE_KINFLATE_EXPECTED_DATA"
	updateKinflateExpected := os.Getenv(updateEnvVar) == "true"

	buf := bytes.NewBuffer([]byte{})

	cmd := NewCmdInflate(buf, os.Stderr)
	cmd.Flags().Set("filename", input)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	actualBytes := buf.Bytes()
	if !updateKinflateExpected {
		expectedBytes, err := ioutil.ReadFile(expected)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(actualBytes, expectedBytes) {
			t.Errorf("%s\ndoesn't equal expected:\n%s\n", actualBytes, expectedBytes)
		}
	} else {
		ioutil.WriteFile(expected, actualBytes, 0644)
	}
}
