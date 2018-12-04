package composer

import (
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/router"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation"
)

type composer struct {
	network string
	port    string
}

type Composer interface {
	Run(doneCh <-chan struct{})
}

func New(network, port string) Composer {
	return &composer{network, port}
}

func (composer *composer) Run(done <-chan struct{}) {
	swapRequests := make(chan foundation.SwapRequest)
	statusUpdates := make(chan foundation.StatusUpdate)
	statusQueries := make(chan foundation.StatusQuery)
	ftSwapRequests := make(chan foundation.SwapRequest)
	ftStatusUpdates := make(chan foundation.StatusUpdate)
	statuses := make(chan foundation.SwapStatus)
	results := make(chan foundation.SwapResult)

	manager, err := keystore.FundManager(composer.network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore()
	if err != nil {
		panic(err)
	}

	passwordHash, err := keystore.LoadPasswordHash(composer.network)
	if err != nil {
		panic(err)
	}

	logger := logger.NewStdOut()

	co.ParBegin(
		func() {
			httpServer := server.NewHttpServer(manager, logger, passwordHash, composer.port)
			httpServer.Run(done, swapRequests, statusQueries)
		},
		func() {
			router := router.New(db.New(ldb), logger)
			router.Run(done, swapRequests, statusUpdates, ftSwapRequests, ftStatusUpdates, statuses)
		},
		func() {
			swapper := swapper.New(callback.New(), binder.NewBuilder(manager, logger), logger)
			swapper.Run(done, ftSwapRequests, results, statusUpdates)
		},
		func() {
			statusHandler := status.New()
			statusHandler.Run(done, statuses, ftStatusUpdates, statusQueries)
		},
	)
}
