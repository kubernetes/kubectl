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

package apply

import (
	"os"
	"testing"
)

const CONFIG_FILE = "config.yaml"

func TestNewCmdApply(t *testing.T) {
	cmd := NewCmdApply(os.Stdin, os.Stdout, os.Stderr)
	if cmd == nil {
		t.Error("Cmd is nil\n")
	}
	// Check the command has the flags
	if !cmd.Flags().HasFlags() {
		t.Error("Cmd missing flags\n")
	}
}

func TestValidateApplyCommand(t *testing.T) {
	options := ApplyOptions{ConfigFile: CONFIG_FILE}
	cmd := NewCmdApply(os.Stdin, os.Stdout, os.Stderr)
	// Validate that nil cmd throws error
	err := options.Validate(nil, []string{})
	if err == nil {
		t.Errorf("Expected error not thrown: missing args")
	}
	// Validate non-empty args throws error
	err = options.Validate(nil, []string{""})
	if err == nil {
		t.Errorf("Expected error not thrown: unexpected args")
	}
	// Validate that correct parameters doesn't throw error
	err = options.Validate(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRunApplyCommand(t *testing.T) {
	options := ApplyOptions{ConfigFile: CONFIG_FILE}
	cmd := NewCmdApply(os.Stdin, os.Stdout, os.Stderr)
	err := options.RunApply(cmd, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
