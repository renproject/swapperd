package swapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSwapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swapper Suite")
}
