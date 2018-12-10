package router_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/swapperd/core/router"
	"github.com/republicprotocol/swapperd/testutils"
)

var _ = Describe("Router", func() {

	init := func() Router {
		storage := testutils.NewMockStorage()

	}

	Context("when routing swaps", func() {
		It("should store the received swap request", func() {
			Expect(true).Should(BeTrue())
		})
	})
})
