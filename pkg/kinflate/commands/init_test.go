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
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdInit(buf, os.Stderr)
	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if buf.String() != "Hello I am init.\n" {
		t.Errorf("unexpected output: %v", buf.String())
	}
}
