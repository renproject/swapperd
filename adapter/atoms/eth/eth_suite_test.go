package eth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eth Suite")
}
