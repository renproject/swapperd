package swapperd

import (
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/router"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
)

func Run(doneCh <-chan struct{}, network, port string) {
	manager, err := keystore.FundManager(network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore()
	if err != nil {
		panic(err)
	}

	passwordHash, err := keystore.LoadPasswordHash(network)
	if err != nil {
		panic(err)
	}

	logger := logger.NewStdOut()

	router := router.New(
		swapper.New(callback.New(), binder.NewBuilder(manager, logger), logger),
		status.New(),
		db.New(ldb),
		logger,
		server.NewHttpServer(manager, logger, passwordHash, port),
	)

	go router.Run(doneCh)
}
