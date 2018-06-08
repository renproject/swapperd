package swap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSwap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swap Suite")
}
