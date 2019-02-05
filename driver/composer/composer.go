package composer

import (
	"github.com/republicprotocol/co-go"
	"github.com/renproject/swapperd/adapter/binder"
	"github.com/renproject/swapperd/adapter/callback"
	"github.com/renproject/swapperd/adapter/db"
	"github.com/renproject/swapperd/adapter/server"
	"github.com/renproject/swapperd/core/wallet"
	"github.com/renproject/swapperd/driver/keystore"
	"github.com/renproject/swapperd/driver/leveldb"
	"github.com/renproject/swapperd/driver/logger"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

const BufferCapacity = 2048

type composer struct {
	server      server.Server
	logger      logrus.FieldLogger
	walletTask  tau.Task
	serviceTask tau.Task
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

	receiver := server.NewReceiver(BufferCapacity)
	serviceTask := server.NewService(BufferCapacity, receiver)
	serviceTask.Send(server.AcceptRequest{})

	server := server.NewHttpServer(BufferCapacity, port, receiver, storage, bc, logger)
	walletTask := wallet.New(BufferCapacity, storage, bc, binder.NewBuilder(bc, logger), callback.New())
	return &composer{server, logger, walletTask, serviceTask}
}

func (composer *composer) Run(done <-chan struct{}) {
	co.ParBegin(
		func() {
			tau.New(tau.NewIO(BufferCapacity), tau.ReduceFunc(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case server.AcceptedRequest:
					composer.walletTask.Send(msg.Message)
					composer.serviceTask.Send(server.AcceptRequest{})
				case tau.Error:
					composer.logger.Error(msg)
				default:
					composer.logger.Errorf("Unexpected message type: %T in compser", msg)
				}
				return nil
			}), composer.walletTask, composer.serviceTask).Run(done)
		},
		func() {
			composer.server.Run(done)
		},
	)
}
