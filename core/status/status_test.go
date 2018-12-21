package status_test

import (
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/swapperd/testutils"
)

var Random *rand.Rand

func init() {
	Random = rand.New(rand.NewSource(time.Now().Unix()))
}

var _ = Describe("Statuses", func() {

	init := func() (Statuses, chan struct{}, chan swap.SwapReceipt, chan swap.ReceiptUpdate, chan ReceiptQuery) {
		logger := logger.NewStdOut()
		store := testutils.NewMockStorage()
		statuses := New(store, logger)
		done := make(chan struct{})
		receipts := make(chan swap.SwapReceipt)
		updates := make(chan swap.ReceiptUpdate)
		queries := make(chan ReceiptQuery)

		return statuses, done, receipts, updates, queries
	}

	Context("when receiving new receipt", func() {
		It("should store the receipt", func() {
			statuses, done, receipts, updates, queries := init()
			defer close(done)
			go statuses.Run(done, receipts, updates, queries)

			test := func(receipt swap.SwapReceipt) bool {
				receipts <- receipt
				responder := make(chan map[swap.SwapID]swap.SwapReceipt, 1)
				query := ReceiptQuery{
					Responder: responder,
				}
				queries <- query

				response := <-responder
				return reflect.DeepEqual(response[receipt.ID], receipt)
			}

			Expect(quick.Check(test, nil)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving new update", func() {
		It("should update the receipt", func() {
			statuses, done, receipts, updates, queries := init()
			defer close(done)
			go statuses.Run(done, receipts, updates, queries)

			test := func(receipt swap.SwapReceipt) bool {
				receipts <- receipt

				receipt.Status = swap.Audited
				update := swap.NewReceiptUpdate(receipt.ID, func(re *swap.SwapReceipt) {
					re.ID = receipt.ID
					re.Status = swap.Audited
				})

				updates <- update

				responder := make(chan map[swap.SwapID]swap.SwapReceipt, 1)
				query := ReceiptQuery{
					Responder: responder,
				}
				queries <- query
				response := <-responder

				return reflect.DeepEqual(response[receipt.ID], receipt)
			}

			Expect(quick.Check(test, nil)).ShouldNot(HaveOccurred())
		})
	})

	Context("when closing one of the input channel", func() {
		It("should stop the statues from running ", func() {
			statuses, done, receipts, updates, queries := init()
			go statuses.Run(done, receipts, updates, queries)

			close(receipts)

			// Expect(receipts).ShouldNot(BeSent(swap.SwapReceipt{}))
			Expect(updates).ShouldNot(BeSent(swap.ReceiptUpdate{}))
			Expect(queries).ShouldNot(BeSent(ReceiptQuery{}))
		})

		It("should stop the statues from running ", func() {
			statuses, done, receipts, updates, queries := init()
			go statuses.Run(done, receipts, updates, queries)

			close(updates)
			time.Sleep(1 * time.Millisecond)

			Expect(receipts).ShouldNot(BeSent(swap.SwapReceipt{}))
			Expect(queries).ShouldNot(BeSent(ReceiptQuery{}))
		})

		It("should stop the statues from running ", func() {
			statuses, done, receipts, updates, queries := init()
			go statuses.Run(done, receipts, updates, queries)

			close(queries)

			Expect(receipts).ShouldNot(BeSent(swap.SwapReceipt{}))
			Expect(updates).ShouldNot(BeSent(swap.ReceiptUpdate{}))
		})
	})
})
