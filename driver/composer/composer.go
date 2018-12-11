package composer

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation/swap"
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
	receipts := make(chan swap.SwapReceipt)

	statusUpdates := make(chan swap.StatusUpdate)
	receiptQueries := make(chan swap.ReceiptQuery)
	balanceQueries := make(chan balance.BalanceQuery)

	wallet, err := keystore.Wallet(composer.homeDir, composer.network)
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

	co.ParBegin(
		func() {
			httpServer := server.NewHttpServer(wallet, storage, logger, passwordHash, composer.port)
			httpServer.Run(done, swaps, receipts, receiptQueries, balanceQueries)
		},
		func() {
			swapper := swapper.New(callback.New(), binder.NewBuilder(wallet, logger), storage, logger)
			swapper.Run(done, swaps, statusUpdates)
		},
		func() {
			statuses := status.New()
			statuses.Run(done, receipts, statusUpdates, receiptQueries)
		},
		func() {
			updateFrequency := 15 * time.Second
			balanceHandler := balance.New(updateFrequency, wallet, logger)
			balanceHandler.Run(done, balanceQueries)
		},
	)
}
