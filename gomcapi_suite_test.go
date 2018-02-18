package mcapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGomcapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomcapi Suite")
}
