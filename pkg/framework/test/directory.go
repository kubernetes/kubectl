package test

import (
	"io/ioutil"
	"os"
)

// Directory holds a path to a directory and knows how to tear down / cleanup that directory
type Directory struct {
	Path    string
	Cleanup func() error
}

func newDirectory() (*Directory, error) {
	path, err := ioutil.TempDir("", "k8s_test_framework_")
	if err != nil {
		return nil, err
	}

	return &Directory{
		Path: path,
		Cleanup: func() error {
			return os.RemoveAll(path)
		},
	}, nil
}
