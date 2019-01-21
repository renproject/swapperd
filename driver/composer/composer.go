package composer

import (
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/core/swapper/delayed"
	"github.com/republicprotocol/swapperd/core/swapper/immediate"
	"github.com/republicprotocol/swapperd/core/swapper/status"
	"github.com/republicprotocol/swapperd/core/transfer"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
)

const BufferCapacity = 128

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
	blockchain, err := keystore.Wallet(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore(composer.homeDir, composer.network)
	if err != nil {
		panic(err)
	}

	storage := db.New(ldb)
	logger := logger.NewStdOut()

	delayedSwapperTask := delayed.New(BufferCapacity, callback.New())
	immediateSwapperTask := immediate.New(BufferCapacity, binder.NewBuilder(blockchain, logger))
	swapStatusTask := status.New(BufferCapacity)

	swapperTask := swapper.New(BufferCapacity, storage, delayedSwapperTask, immediateSwapperTask, swapStatusTask)
	walletTask := transfer.New(BufferCapacity, blockchain, storage)

	httpServer := server.NewHttpServer(blockchain, logger, swapperTask, walletTask, composer.port)
	httpServer.Run(done)
}
