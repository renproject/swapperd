package balance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBalance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Balance Suite")
}
