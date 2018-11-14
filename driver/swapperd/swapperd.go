package swapperd

import (
	"fmt"
	"net"
	"net/http"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/callback"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/adapter/storage"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation"
)

func Run(doneCh <-chan struct{}, network, port string) {
	swaps := make(chan swapper.Swap)
	statuses := make(chan foundation.SwapStatus)
	statusQueries := make(chan status.Query)

	manager, err := keystore.FundManager(network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore()
	if err != nil {
		panic(err)
	}

	go co.ParBegin(
		func() {
			authenticator, err := keystore.LoadAuthenticator(network)
			if err != nil {
				panic(err)
			}
			handler := server.NewHandler(authenticator, manager, swaps, statusQueries)
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
			if err != nil {
				panic(err)
			}
			go func() {
				if err := http.Serve(listener, handler); err != nil {
					panic(err)
				}
			}()
			<-doneCh
			listener.Close()
		},
		func() {
			stdLogger := logger.NewStdOut()
			builder := binder.NewBuilder(manager, stdLogger)
			storage := storage.New(ldb)
			callback := callback.New()
			swapper := swapper.New(callback, builder, storage, stdLogger)
			swapper.Run(doneCh, swaps, statuses)
		},
		func() {
			monitor := status.New()
			monitor.Run(doneCh, statuses, statusQueries)
		},
	)
}
