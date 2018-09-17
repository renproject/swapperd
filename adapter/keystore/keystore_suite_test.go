package keystore_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKeystore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keystore Suite")
}
