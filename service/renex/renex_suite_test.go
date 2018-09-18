package renex_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRenex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Renex Suite")
}
