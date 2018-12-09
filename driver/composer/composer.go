package composer

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/core/router"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation"
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
	swapRequests := make(chan foundation.SwapRequest)
	statusUpdates := make(chan foundation.StatusUpdate)
	statusQueries := make(chan foundation.StatusQuery)
	balanceQueries := make(chan balance.BalanceQuery)
	ftSwapRequests := make(chan foundation.SwapRequest)
	ftStatusUpdates := make(chan foundation.StatusUpdate)
	receipts := make(chan foundation.SwapStatus)
	results := make(chan foundation.SwapResult)

	wallet, err := keystore.Wallet(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	passwordHash, err := keystore.LoadPasswordHash(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	logger := logger.NewStdOut()

	co.ParBegin(
		func() {
			httpServer := server.NewHttpServer(wallet, logger, passwordHash, composer.port)
			httpServer.Run(done, swapRequests, statusQueries, balanceQueries)
		},
		func() {
			router := router.New(db.New(ldb), logger)
			router.Run(done, swapRequests, statusUpdates, ftSwapRequests, ftStatusUpdates, receipts)
		},
		func() {
			swapper := swapper.New(callback.New(), binder.NewBuilder(wallet, logger), logger)
			swapper.Run(done, ftSwapRequests, results, statusUpdates)
		},
		func() {
			statuses := status.New()
			statuses.Run(done, receipts, ftStatusUpdates, statusQueries)
		},
		func() {
			updateFrequency := 15 * time.Second
			balanceHandler := balance.New(updateFrequency, wallet, logger)
			balanceHandler.Run(done, balanceQueries)
		},
	)
}
