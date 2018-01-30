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
)

var _ FileSystem = OSFS{}

// OSFS implements FileSystem using the local filesystem.
type OSFS struct{}

// Create creates a file given the filename.
func (OSFS) Create(name string) (File, error) { return os.Create(name) }

// Open opens a file given the filename.
func (OSFS) Open(name string) (File, error) { return os.Open(name) }

// Stat return an interface which has all the information regarding the file.
func (OSFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
