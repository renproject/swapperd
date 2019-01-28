package status_test

import (
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	"github.com/republicprotocol/swapperd/testutils"

	"github.com/republicprotocol/tau"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/core/wallet/swapper/status"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

var Random *rand.Rand

func init() {
	Random = rand.New(rand.NewSource(time.Now().Unix()))
}

var _ = Describe("Status Task", func() {

	init := func() (tau.Task, chan struct{}) {
		return New(testutils.DefaultQuickCheckConfig.MaxCount), make(chan struct{})
	}

	Context("when receiving new receipt", func() {
		It("should store the receipt", func() {
			statusTask, done := init()
			defer close(done)
			go statusTask.Run(done)

			test := func(receipt Receipt) bool {
				statusTask.IO().InputWriter() <- receipt
				responder := make(chan map[swap.SwapID]swap.SwapReceipt, 1)
				query := ReceiptQuery{
					Responder: responder,
				}
				statusTask.IO().InputWriter() <- query
				response := <-responder
				return reflect.DeepEqual(response[receipt.ID], swap.SwapReceipt(receipt))
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving new update", func() {
		It("should update the receipt", func() {
			statusTask, done := init()
			defer close(done)
			go statusTask.Run(done)

			test := func(receipt Receipt) bool {
				statusTask.IO().InputWriter() <- receipt

				receipt.Status = swap.Audited
				update := ReceiptUpdate(swap.NewReceiptUpdate(receipt.ID, func(re *swap.SwapReceipt) {
					re.ID = receipt.ID
					re.Status = swap.Audited
				}))

				statusTask.IO().InputWriter() <- update

				responder := make(chan map[swap.SwapID]swap.SwapReceipt, 1)
				query := ReceiptQuery{
					Responder: responder,
				}
				statusTask.IO().InputWriter() <- query
				response := <-responder

				return reflect.DeepEqual(response[receipt.ID], swap.SwapReceipt(receipt))
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving an unknown message type", func() {
		It("should return an error", func() {
			statusTask, done := init()
			defer close(done)
			go statusTask.Run(done)

			statusTask.IO().InputWriter() <- tau.RandomMessage{}
			err := <-statusTask.IO().OutputReader()
			_, ok := err.(tau.Error)
			Expect(ok).Should(BeTrue())
		})
	})
})
