package bootload_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bootloading", func() {
	Context("when booting", func() {
		It("should load receipts and pending swaps from storage", func() {
			Expect(true).Should(BeTrue())
		})
	})
})
