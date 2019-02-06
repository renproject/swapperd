package immediate_test

import (
	"testing/quick"

	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/swapperd/testutils"
	"github.com/republicprotocol/tau"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/swapperd/core/wallet/swapper/immediate"
)

var _ = Describe("Immediate Swap Task", func() {

	init := func() (tau.Task, chan struct{}) {
		contractBuilder := NewMockContractBuilder()
		return New(testutils.DefaultQuickCheckConfig.MaxCount, contractBuilder), make(chan struct{})
	}

	handleSwapResponse := func(msg tau.Message, blob swap.SwapBlob) bool {
		timeLock := uint64(blob.TimeLock)
		if timeLock%36 == 8 {
			_, ok := msg.(tau.Error)
			return ok
		}

		messages := msg.(tau.MessageBatch)
		update := messages[0].(ReceiptUpdate)
		receipt := swap.NewSwapReceipt(blob)
		update.Update(&receipt)

		switch timeLock % 9 {
		case 0:
			_, ok := messages[1].(tau.Error)
			if blob.ShouldInitiateFirst {
				return ok && len(messages) == 2 && receipt.Status == swap.Inactive
			}
			return ok && len(messages) == 2 && receipt.Status == swap.Audited
		case 1:
			if !blob.ShouldInitiateFirst {
				_, msg1ok := messages[1].(tau.Error)
				deleteSwap, msg2ok := messages[2].(DeleteSwap)
				return msg1ok && msg2ok && len(messages) == 3 && receipt.Status == swap.Expired && deleteSwap.ID == blob.ID
			}
			if timeLock%18 == 1 {
				_, ok := messages[1].(tau.Error)
				return ok && len(messages) == 2 && receipt.Status == swap.RefundFailed
			}
			deleteSwap, ok := messages[1].(DeleteSwap)
			return ok && len(messages) == 2 && deleteSwap.ID == blob.ID && receipt.Status == swap.Refunded
		case 2:
			return len(messages) == 1
		case 3:
			_, ok := messages[1].(tau.Error)
			return ok && len(messages) == 2 && receipt.Status == swap.AuditPending
		case 4:
			_, ok := messages[1].(tau.Error)
			if blob.ShouldInitiateFirst {
				return ok && len(messages) == 2 && receipt.Status == swap.Audited
			}
			return ok && len(messages) == 2 && receipt.Status == swap.AuditedSecret
		case 5:
			_, ok := messages[1].(tau.Error)
			return blob.ShouldInitiateFirst || ok && len(messages) == 2 && receipt.Status == swap.Initiated
		case 6:
			return blob.ShouldInitiateFirst || len(messages) == 1
		case 7:
			if timeLock%18 == 7 {
				_, ok := messages[1].(tau.Error)
				return blob.ShouldInitiateFirst || ok && len(messages) == 2 && receipt.Status == swap.RefundFailed
			}
			deleteSwap, ok := messages[1].(DeleteSwap)
			return blob.ShouldInitiateFirst || ok && len(messages) == 2 && deleteSwap.ID == blob.ID && receipt.Status == swap.Refunded
		case 8:
			deleteSwap, ok := messages[1].(DeleteSwap)
			return ok && len(messages) == 2 && receipt.Status == swap.Redeemed && deleteSwap.ID == blob.ID
		default:
			return false
		}
	}

	Context("when receiving new immediate swap request", func() {
		It("should return receipt update and new swap on nil error", func() {
			immediateTask, done := init()
			defer close(done)
			go immediateTask.Run(done)

			test := func(blob swap.SwapBlob) bool {
				request := NewSwapRequest(blob, blockchain.Cost{}, blockchain.Cost{})
				immediateTask.IO().InputWriter() <- request
				response := <-immediateTask.IO().OutputReader()

				return handleSwapResponse(response, blob)
			}

			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})

		Context("when receiveng new tick", func() {
			It("should return ", func() {
				immediateTask, done := init()
				defer close(done)
				go immediateTask.Run(done)

				test := func(blob swap.SwapBlob) bool {
					request := NewSwapRequest(blob, blockchain.Cost{}, blockchain.Cost{})
					immediateTask.IO().InputWriter() <- request
					<-immediateTask.IO().OutputReader()
					immediateTask.IO().InputWriter() <- tau.Tick{}
					<-immediateTask.IO().OutputReader()
					return true
				}

				Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
			})
		})

		Context("when receiving an unknown message type", func() {
			It("should return an error", func() {
				immediateTask, done := init()
				defer close(done)
				go immediateTask.Run(done)

				test := func() bool {
					immediateTask.IO().InputWriter() <- tau.RandomMessage{}
					err := <-immediateTask.IO().OutputReader()
					_, ok := err.(tau.Error)
					return ok
				}

				Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
			})
		})
	})
})
