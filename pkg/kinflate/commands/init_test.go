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
	"path"
	"testing"

	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

func TestInitHappyPath(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	fakeFS := fs.MakeFakeFS()
	cmd := NewCmdInit(buf, os.Stderr, fakeFS)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := fakeFS.Open(path.Join(appname, "Kube-manifest.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file := f.(*fs.FakeFile)
	if !file.ContentMatches([]byte(manifestTemplate)) {
		t.Fatalf("actual: %v doesn't match expected: %v",
			string(file.GetContent()), manifestTemplate)
	}
	f, err = fakeFS.Open(path.Join(appname, "deployment.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file = f.(*fs.FakeFile)
	if !file.ContentMatches([]byte(deploymentTemplate)) {
		t.Fatalf("actual: %v doesn't match expected: %v",
			string(file.GetContent()), deploymentTemplate)
	}
	f, err = fakeFS.Open(path.Join(appname, "service.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file = f.(*fs.FakeFile)
	if !file.ContentMatches([]byte(serviceTemplate)) {
		t.Fatalf("actual: %v doesn't match expected: %v",
			string(file.GetContent()), serviceTemplate)
	}
}

func TestInitFileAlreadyExist(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	fakeFS := fs.MakeFakeFS()
	fakeFS.Mkdir(appname, 0766)

	cmd := NewCmdInit(buf, os.Stderr, fakeFS)
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != `"helloworld" already exists` {
		t.Fatalf("actual err: %v doesn't match expected error: %v", err, `"helloworld" already exists`)
	}
}
