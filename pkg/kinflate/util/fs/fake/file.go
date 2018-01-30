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

package fs

import (
	"os"

	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

var _ fs.File = &FakeFile{}

// FakeFile implements FileSystem using a fake in-memory filesystem.
type FakeFile struct {
	content []byte
	open    bool
}

// Close closes a file.
func (f *FakeFile) Close() error {
	f.open = false
	return nil
}

// Read reads a file's content.
func (f *FakeFile) Read(p []byte) (n int, err error) {
	return len(p), nil
}

// Write writes bytes to a file
func (f *FakeFile) Write(p []byte) (n int, err error) {
	f.content = p
	return len(p), nil
}

// Stat returns an interface which has all the information regarding the file.
func (f *FakeFile) Stat() (os.FileInfo, error) {
	return nil, nil
}
