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
	"testing"
)

// TODO: real tests
// e.g. make an inmemory file system, put yaml in there, inflate it
// to a buffer, compare to expected results, etc.
// a script in there, have script write file
func TestTrueMain(t *testing.T) {
	exrr := TestableMain()
	if exrr != nil {
		t.Errorf("Unexpected error: %v", exrr)
	}
}
