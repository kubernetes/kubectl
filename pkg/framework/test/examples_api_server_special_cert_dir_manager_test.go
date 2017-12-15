package test_test

import (
	"fmt"

	"io/ioutil"
	"os"

	. "k8s.io/kubectl/pkg/framework/test"
)

func ExampleAPIServer_credHubCertDirManager() {
	apiServer := &APIServer{
		CertDirManager: NewCredHubCertDirManager(),
	}
	fmt.Println(apiServer)
}

func credHubLoader(dir string, prefix string) (string, error) {
	tempDir, _ := ioutil.TempDir("/var/cache/kube/cert", prefix+"cred-hub-")
	loadCertsFromCredHub(tempDir)
	return tempDir, nil
}

func credHubSaver(tempDir string) error {
	saveCertsToCredHub(tempDir)
	return os.RemoveAll(tempDir)
}

func NewCredHubCertDirManager() *TempDirManager {
	return &TempDirManager{
		Maker:   credHubLoader,
		Remover: credHubSaver,
	}
}

func loadCertsFromCredHub(dir string) { /* to de implemented */ }
func saveCertsToCredHub(dir string)   { /* to be implemented */ }
