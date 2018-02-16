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
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"k8s.io/apimachinery/pkg/util/sets"
)

type InflateTestCase struct {
	Description string   `yaml:"description"`
	Args        []string `yaml:"args"`
	Filename    string   `yaml:"filename"`
	// path to the file that contains the expected output
	ExpectedStdout string `yaml:"expectedStdout"`
}

func TestInflate(t *testing.T) {
	const updateEnvVar = "UPDATE_KINFLATE_EXPECTED_DATA"
	updateKinflateExpected := os.Getenv(updateEnvVar) == "true"

	var (
		name     string
		testcase InflateTestCase
	)

	testcases := sets.NewString()
	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "testdata" {
			return nil
		}
		name := filepath.Base(path)
		if info.IsDir() {
			if strings.HasPrefix(name, "testcase-") {
				testcases.Insert(strings.TrimPrefix(name, "testcase-"))
			}
			return filepath.SkipDir
		}
		return nil
	})
	// sanity check that we found the right folder
	if !testcases.Has("simple") {
		t.Fatalf("Error locating kinflate inflate testcases")
	}

	for _, testcaseName := range testcases.List() {
		t.Run(testcaseName, func(t *testing.T) {
			name = testcaseName
			testcase = InflateTestCase{}
			testcaseDir := filepath.Join("testdata", "testcase-"+name)
			testcaseData, err := ioutil.ReadFile(filepath.Join(testcaseDir, "test.yaml"))
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if err := yaml.Unmarshal(testcaseData, &testcase); err != nil {
				t.Fatalf("%s: %v", name, err)
			}

			buf := bytes.NewBuffer([]byte{})

			cmd := newCmdInflate(buf, os.Stderr)
			cmd.Flags().Set("filename", testcase.Filename)

			err = cmd.Execute()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			actualBytes := buf.Bytes()
			if !updateKinflateExpected {
				expectedBytes, err := ioutil.ReadFile(testcase.ExpectedStdout)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(actualBytes, expectedBytes) {
					t.Errorf("%s\ndoesn't equal expected:\n%s\n", actualBytes, expectedBytes)
				}
			} else {
				ioutil.WriteFile(testcase.ExpectedStdout, actualBytes, 0644)
			}

		})
	}

}
