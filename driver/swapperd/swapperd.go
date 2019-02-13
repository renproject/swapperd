package swapperd

import (
	"github.com/renproject/swapperd/adapter/binder"
	"github.com/renproject/swapperd/adapter/callback"
	"github.com/renproject/swapperd/adapter/db"
	"github.com/renproject/swapperd/adapter/server"
	"github.com/renproject/swapperd/core/wallet"
	"github.com/renproject/swapperd/driver/keystore"
	"github.com/renproject/swapperd/driver/leveldb"
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

const BufferCapacity = 2048

type swapperd struct {
	server      server.Server
	logger      logrus.FieldLogger
	walletTask  tau.Task
	serviceTask tau.Task
}

type Swapperd interface {
	Run(doneCh <-chan struct{})
}

func New(version, homeDir, network, port string, logger logrus.FieldLogger) Swapperd {
	ldb, err := leveldb.NewStore(homeDir, network)
	if err != nil {
		panic(err)
	}

	storage := db.New(ldb)
	bc, err := keystore.Wallet(homeDir, network)
	if err != nil {
		panic(err)
	}

	receiver := server.NewReceiver(BufferCapacity)
	serviceTask := server.NewService(BufferCapacity, receiver)
	serviceTask.Send(server.AcceptRequest{})

	server := server.NewHttpServer(BufferCapacity, port, version, receiver, storage, bc, logger)
	walletTask := wallet.New(BufferCapacity, storage, bc, binder.NewBuilder(bc, logger), callback.New())
	return &swapperd{server, logger, walletTask, serviceTask}
}

func (swapperd *swapperd) Run(done <-chan struct{}) {
	co.ParBegin(
		func() {
			tau.New(tau.NewIO(BufferCapacity), tau.ReduceFunc(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case server.AcceptedRequest:
					swapperd.walletTask.Send(msg.Message)
					swapperd.serviceTask.Send(server.AcceptRequest{})
				case tau.Error:
					swapperd.logger.Error(msg)
				default:
					swapperd.logger.Errorf("Unexpected message type: %T in compser", msg)
				}
				return nil
			}), swapperd.walletTask, swapperd.serviceTask).Run(done)
		},
		func() {
			swapperd.server.Run(done)
		},
	)
}
