package test

import (
	"io/ioutil"
	"os"
)

// CleanableDirectory holds a path to a directory and knows how to tear down / cleanup that directory
type CleanableDirectory struct {
	Path    string
	Cleanup func() error
}

func newDirectory() (*CleanableDirectory, error) {
	path, err := ioutil.TempDir("", "k8s_test_framework_")
	if err != nil {
		return nil, err
	}

	return &CleanableDirectory{
		Path: path,
		Cleanup: func() error {
			return os.RemoveAll(path)
		},
	}, nil
}
