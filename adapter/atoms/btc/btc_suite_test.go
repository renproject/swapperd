package btc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBtc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Btc Suite")
}
