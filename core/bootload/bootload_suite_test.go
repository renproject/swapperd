package bootload_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBootload(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bootload Suite")
}
