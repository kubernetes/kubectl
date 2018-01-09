package pluginutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFramework(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Framework Suite")
}
