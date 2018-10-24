package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/balance"
	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/adapter/storage"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/leveldb"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation"
)

func main() {
	network := flag.String("network", "testnet", "Which network to use")
	port := flag.Int64("port", 7777, "Which network to use")
	flag.Parse()

	swaps := make(chan swapper.Swap)
	statuses := make(chan foundation.SwapStatus)
	statusQueries := make(chan status.Query)
	balanceQueries := make(chan balance.Query)

	accounts, err := keystore.LoadAccounts(*network)
	if err != nil {
		panic(err)
	}

	authenticator, err := keystore.LoadAuthenticator(*network)
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore()
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go co.ParBegin(
		func() {
			handler := server.NewHandler(authenticator, swaps, statusQueries, balanceQueries)
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
			if err != nil {
				panic(err)
			}
			go func() {
				if err := http.Serve(listener, handler); err != nil {
					panic(err)
				}
			}()
			<-done
			listener.Close()
		},
		func() {
			stdLogger := logger.NewStdOut()
			builder := binder.NewBuilder(accounts, stdLogger)
			storage := storage.New(ldb)
			swapper := swapper.New(builder, storage, stdLogger)
			swapper.Run(done, swaps, statuses)
		},
		func() {
			monitor := status.New()
			monitor.Run(done, statuses, statusQueries)
		},
		func() {
			balanceBook := balance.New(accounts)
			balanceBook.Run(done, balanceQueries)
		},
	)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}