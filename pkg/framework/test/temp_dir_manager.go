package test

import (
	"io/ioutil"
	"os"
)

// TempDirMaker can create directories.
type TempDirMaker func(dir, prefix string) (name string, err error)

// TempDirRemover can delete directories
type TempDirRemover func(dir string) error

// NewTempDirManager returns a new manager for creation and deleteion of temporary directories.
func NewTempDirManager() *TempDirManager {
	return &TempDirManager{
		Maker:   ioutil.TempDir,
		Remover: os.RemoveAll,
	}
}

// TempDirManager knows when to call the directory maker and remover and keeps track of created directories.
type TempDirManager struct {
	Maker   TempDirMaker
	Remover TempDirRemover
	dir     string
}

// Create knows how to create a temporary directory and how to keep track of it.
func (t *TempDirManager) Create() (string, error) {
	if t.dir == "" {
		dir, err := t.Maker("", "kube-test-framework-")
		if err != nil {
			return "", err
		}
		t.dir = dir
	}
	return t.dir, nil
}

// Destroy knows how to destroy a previously created directory.
func (t *TempDirManager) Destroy() error {
	if t.dir != "" {
		err := t.Remover(t.dir)
		t.dir = ""
		return err
	}
	return nil
}
