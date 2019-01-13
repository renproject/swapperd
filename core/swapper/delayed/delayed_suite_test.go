package delayed_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDelayed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delayed Suite")
}
