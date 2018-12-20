package composer

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/delayed"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/core/wallet"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
)

type composer struct {
	homeDir string
	network string
	port    string
}

type Composer interface {
	Run(doneCh <-chan struct{})
}

func New(homeDir, network, port string) Composer {
	return &composer{homeDir, network, port}
}

func (composer *composer) Run(done <-chan struct{}) {
	swaps := make(chan swap.SwapBlob)
	delayedSwaps := make(chan swap.SwapBlob)
	receipts := make(chan swap.SwapReceipt)
	receiptUpdates := make(chan swap.ReceiptUpdate)
	receiptQueries := make(chan status.ReceiptQuery)

	blockchain, err := keystore.Wallet(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	storage := db.New(ldb)

	passwordHash, err := keystore.LoadPasswordHash(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	logger := logger.NewStdOut()
	walletTask := wallet.New(128, blockchain, storage, logger)

	co.ParBegin(
		func() {
			httpServer := server.NewHttpServer(blockchain, storage, logger, passwordHash, composer.port)
			httpServer.Run(done, swaps, delayedSwaps, receipts, receiptQueries, walletTask.IO())
		},
		func() {
			delayedCallback := delayed.New(callback.New(), storage, logger)
			delayedCallback.Run(done, delayedSwaps, swaps, receiptUpdates)
		},
		func() {
			swapper := swapper.New(binder.NewBuilder(blockchain, logger), storage, logger)
			swapper.Run(done, swaps, receiptUpdates)
		},
		func() {
			statuses := status.New(storage, logger)
			statuses.Run(done, receipts, receiptUpdates, receiptQueries)
		},
		func() {
			ticker := time.NewTicker(30 * time.Second)
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					walletTask.IO().InputWriter() <- tau.NewTick(time.Now())
				}
			}
		},
		func() {
			walletTask.Run(done)
		},
	)
}
