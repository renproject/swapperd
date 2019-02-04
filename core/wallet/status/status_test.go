package status_test

import (
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	"github.com/renproject/swapperd/testutils"

	"github.com/republicprotocol/tau"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/swapperd/core/wallet/status"

	"github.com/renproject/swapperd/foundation/swap"
)

var Random *rand.Rand

func init() {
	Random = rand.New(rand.NewSource(time.Now().Unix()))
}

var _ = Describe("Status Task", func() {

	init := func() (*MockStorage, tau.Reducer) {
		storage := NewMockStorage()
		return storage, NewReducer(storage)
	}

	Context("when receiving new receipt", func() {
		It("should store the receipt", func() {
			storage, reducer := init()

			test := func(receipt Receipt) bool {
				Expect(reducer.Reduce(receipt)).Should(BeNil())
				statusReceipt, err := storage.Receipt(receipt.ID)
				Expect(err).Should(BeNil())
				return reflect.DeepEqual(statusReceipt, swap.SwapReceipt(receipt))
			}
			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving new update", func() {
		It("should update the receipt", func() {
			storage, reducer := init()
			test := func(receipt Receipt) bool {
				Expect(reducer.Reduce(receipt)).Should(BeNil())

				receipt.Status = swap.Audited
				update := ReceiptUpdate(swap.NewReceiptUpdate(receipt.ID, func(re *swap.SwapReceipt) {
					re.ID = receipt.ID
					re.Status = swap.Audited
				}))
				Expect(reducer.Reduce(update)).Should(BeNil())

				statusReceipt, err := storage.Receipt(receipt.ID)
				Expect(err).Should(BeNil())
				return reflect.DeepEqual(statusReceipt, swap.SwapReceipt(receipt))
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving an unknown message type", func() {
		It("should return an error", func() {
			_, reducer := init()
			msg := reducer.Reduce(tau.RandomMessage{})
			_, ok := msg.(tau.Error)
			Expect(ok).Should(BeTrue())
		})
	})
})
