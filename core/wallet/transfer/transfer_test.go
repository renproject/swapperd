package transfer_test

import (
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/testutils"
	"github.com/republicprotocol/tau"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/swapperd/core/wallet/transfer"
)

const BufferCapacity = 1024

var DefaultQuickCheckConfig = &quick.Config{
	Rand:     rand.New(rand.New(rand.NewSource(time.Now().Unix()))),
	MaxCount: 1024,
}

var _ = Describe("Transfer Task", func() {
	Context("when receiving an unknown message type", func() {
		It("should return an error", func() {
			bc := testutils.NewMockBlockchain(map[blockchain.TokenName]blockchain.Balance{})
			storage := testutils.NewMockStorage()
			transferTask := New(BufferCapacity, bc, storage)
			done := make(chan struct{})
			defer close(done)
			go transferTask.Run(done)
			test := func() bool {
				transferTask.IO().InputWriter() <- tau.RandomMessage{}
				err := <-transferTask.IO().OutputReader()
				_, ok := err.(tau.Error)
				return ok
			}
			Expect(quick.Check(test, testutils.DefaultQuickCheckConfig)).ShouldNot(HaveOccurred())
		})
	})
})
