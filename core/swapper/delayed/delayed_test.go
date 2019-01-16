package delayed_test

import (
	"fmt"
	"reflect"
	"testing/quick"

	"github.com/republicprotocol/swapperd/testutils"

	"github.com/republicprotocol/tau"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/core/swapper/delayed"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

var _ = Describe("Delayed Swap Task", func() {

	init := func(err error) (tau.Task, chan struct{}) {
		callback := testutils.NewMockCallback(err)
		return New(testutils.DefaultQuickCheckConfig.MaxCount, callback), make(chan struct{})
	}

	Context("when receiving new delayed swap request", func() {
		It("should return receipt update and new swap on nil error", func() {
			delayedTask, done := init(nil)
			defer close(done)
			go delayedTask.Run(done)

			test := func(request DelayedSwapRequest) bool {
				delayedTask.IO().InputWriter() <- request
				response := <-delayedTask.IO().OutputReader()

				msg, ok := response.(tau.MessageBatch)
				Expect(ok).Should(BeTrue())

				update, ok := msg[0].(ReceiptUpdate)
				Expect(ok).Should(BeTrue())
				receipt := swap.NewSwapReceipt(swap.SwapBlob(request))
				update.Update(&receipt)

				return reflect.DeepEqual(swap.SwapBlob(request), swap.SwapBlob(msg[1].(SwapRequest)))
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})

		It("should return receipt update and delete swap on ErrSwapCancelled", func() {
			delayedTask, done := init(ErrSwapCancelled)
			defer close(done)
			go delayedTask.Run(done)

			test := func(request DelayedSwapRequest) bool {
				delayedTask.IO().InputWriter() <- request
				response := <-delayedTask.IO().OutputReader()

				msg, ok := response.(tau.MessageBatch)
				Expect(ok).Should(BeTrue())

				update, ok := msg[0].(ReceiptUpdate)
				Expect(ok).Should(BeTrue())
				receipt := swap.NewSwapReceipt(swap.SwapBlob(request))
				update.Update(&receipt)

				return reflect.DeepEqual(DeleteSwap{request.ID}, msg[1].(DeleteSwap))
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})

		It("should not return anything on ErrSwapDetailsUnavailable", func() {
			delayedTask, done := init(ErrSwapDetailsUnavailable)
			defer close(done)
			go delayedTask.Run(done)

			test := func(request DelayedSwapRequest) bool {
				delayedTask.IO().InputWriter() <- request
				return len(delayedTask.IO().OutputReader()) == 0
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})

		It("should return tau error when an unknown error is returned", func() {
			test := func(request DelayedSwapRequest, errString string) bool {
				delayedTask, done := init(fmt.Errorf(errString))
				defer close(done)
				go delayedTask.Run(done)

				delayedTask.IO().InputWriter() <- request
				response := <-delayedTask.IO().OutputReader()
				_, ok := response.(tau.Error)
				return ok
			}
			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiveng new tick", func() {
		It("should return ", func() {
			statusTask, done := init(ErrSwapDetailsUnavailable)
			defer close(done)
			go statusTask.Run(done)
			test := func(request DelayedSwapRequest) bool {
				statusTask.IO().InputWriter() <- request
				statusTask.IO().InputWriter() <- tau.Tick{}
				return true
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})

	Context("when receiving an unknown message type", func() {
		It("should return an error", func() {
			delayedTask, done := init(nil)
			defer close(done)
			go delayedTask.Run(done)

			test := func() bool {
				delayedTask.IO().InputWriter() <- tau.RandomMessage{}
				err := <-delayedTask.IO().OutputReader()
				_, ok := err.(tau.Error)
				return ok
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})
})
