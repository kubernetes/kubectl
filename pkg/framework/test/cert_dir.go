package test

import (
	"io/ioutil"
	"os"
)

// CertDir holds a path to a directory and knows how to tear down / cleanup that directory
type CertDir struct {
	Path    string
	Cleanup func() error
}

func newCertDir() (*CertDir, error) {
	path, err := ioutil.TempDir("", "cert_dir-")
	if err != nil {
		return nil, err
	}

	return &CertDir{
		Path: path,
		Cleanup: func() error {
			return os.RemoveAll(path)
		},
	}, nil
}
