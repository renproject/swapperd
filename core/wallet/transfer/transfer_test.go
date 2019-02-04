package transfer_test

import (
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/swapperd/core/wallet/transfer"
)

const BufferLimit = 1024

var DefaultQuickCheckConfig = &quick.Config{
	Rand:     rand.New(rand.New(rand.NewSource(time.Now().Unix()))),
	MaxCount: 1024,
}

var _ = Describe("Transfer Task", func() {

	responseHandler := func(reducer tau.ReduceFunc) (tau.Task, chan struct{}) {
		storage := NewMockStorage()
		bc := NewMockBlockhain()
		doneCh := make(chan struct{})
		transferTask := New(BufferLimit, bc, storage)
		go tau.New(tau.NewIO(BufferLimit), reducer).Run(doneCh)
		return transferTask, doneCh
	}

	Context("when receiving a bootload request", func() {
		It("should return an error or bootload properly", func() {
			transferTask, doneCh := responseHandler(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case tau.Error:
					return nil
				default:
					Expect(msg).ShouldNot(HaveOccurred())
					return nil
				}
			})
			defer close(doneCh)

			test := func(blob swap.SwapBlob) bool {
				transferTask.IO().InputWriter() <- Bootload{}
				return true
			}
			Expect(quick.Check(test, DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})

		// Context("when receiveng new tick", func() {
		// 	It("should return ", func() {
		// 		immediateTask, done := init()
		// 		defer close(done)
		// 		go immediateTask.Run(done)

		// 		test := func(blob swap.SwapBlob) bool {
		// 			request := NewSwapRequest(blob, blockchain.Cost{}, blockchain.Cost{})
		// 			immediateTask.IO().InputWriter() <- request
		// 			<-immediateTask.IO().OutputReader()
		// 			immediateTask.IO().InputWriter() <- tau.Tick{}
		// 			<-immediateTask.IO().OutputReader()
		// 			return true
		// 		}

		// 		Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		// 	})
		// })

		// Context("when receiving an unknown message type", func() {
		// 	It("should return an error", func() {
		// 		immediateTask, done := init()
		// 		defer close(done)
		// 		go immediateTask.Run(done)

		// 		test := func() bool {
		// 			immediateTask.IO().InputWriter() <- tau.RandomMessage{}
		// 			err := <-immediateTask.IO().OutputReader()
		// 			_, ok := err.(tau.Error)
		// 			return ok
		// 		}

		// 		Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		// 	})
		// })
	})
})
