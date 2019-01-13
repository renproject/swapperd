package swapper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSwapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swapper Suite")
}
