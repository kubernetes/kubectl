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

package main

import (
	"io"
	"os"

	"k8s.io/kubectl/pkg/apply"
)

// TestableMain allows test coverage for main.
func TestableMain(stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := apply.NewCmdApply(stdin, stdout, stderr)
	err := cmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := TestableMain(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
