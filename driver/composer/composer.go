package composer

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/db"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/core/wallet"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

const BufferCapacity = 2048

type composer struct {
	server     server.Server
	logger     logrus.FieldLogger
	walletTask tau.Task
}

type Composer interface {
	Run(doneCh <-chan struct{})
}

func New(homeDir, network, port string) Composer {
	ldb, err := leveldb.NewStore(homeDir, network)
	if err != nil {
		panic(err)
	}

	storage := db.New(ldb)
	logger := logger.NewStdOut()

	bc, err := keystore.Wallet(homeDir, network)
	if err != nil {
		panic(err)
	}
	server := server.NewHttpServer(BufferCapacity, port, bc, logger)

	walletTask := wallet.New(BufferCapacity, storage, bc, binder.NewBuilder(bc, logger), callback.New())
	return &composer{server, logger, walletTask}
}

func (composer *composer) Run(done <-chan struct{}) {
	io := tau.NewIO(BufferCapacity)
	co.ParBegin(
		func() {
			io.InputWriter() <- tau.NewTick(time.Now())
		},
		func() {
			tau.New(io, tau.ReduceFunc(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case tau.Tick:
					composer.walletTask.Send(msg)
				case wallet.AcceptRequest:
					composer.handleAcceptRequest()
				case tau.Error:
					composer.logger.Error(msg)
				default:
					composer.logger.Errorf("Unexpected message type: %T in compser", msg)
				}
				return nil
			}), composer.walletTask).Run(done)
		},
		func() {
			composer.server.Run(done)
		},
	)
}

func (composer *composer) handleAcceptRequest() {
	msg, err := composer.server.Receive()
	if err != nil {
		composer.logger.Error(err)
		return
	}

	switch msg := msg.(type) {
	case tau.Tick, wallet.TransferRequest, wallet.SwapperRequest:
		composer.walletTask.Send(msg)
	default:
		composer.logger.Errorf("Unexpected message type: %T from server", msg)
	}
	return
}
