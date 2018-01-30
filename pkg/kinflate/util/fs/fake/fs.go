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

var _ fs.FileSystem = &FakeFS{}

// FakeFS implements FileSystem using a fake in-memory filesystem.
type FakeFS struct{ m map[string]*FakeFile }

// Create creates a file given the filename.
func (fs *FakeFS) Create(name string) (fs.File, error) {
	if fs.m == nil {
		fs.m = map[string]*FakeFile{}
	}
	fs.m[name] = &FakeFile{}
	return fs.m[name], nil
}

// Open opens a file given the filename.
func (fs *FakeFS) Open(name string) (fs.File, error) {
	if fs.m == nil {
		fs.m = map[string]*FakeFile{}
	}
	fs.m[name] = &FakeFile{open: true}
	return fs.m[name], nil
}

// Stat return an interface which has all the information regarding the file.
func (fs *FakeFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
