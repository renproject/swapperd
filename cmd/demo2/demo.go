package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/republicprotocol/co-go"

	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/adapter/storage"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/driver/store/leveldb"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/republicprotocol/swapperd/utils"
)

// Alice (Has BTC)
// 0x02E596033563647bd3c701622485a8BCF4026175
// mqZ1c5BbVTsaugF2evoxv7HdoEapJ6by5q

// Bob (Has ETH)
// 0xD64fe624faA4c1a858Ba665c50E4f69F1eC85218
// msJSYGB16KqS2yMCC9AA4HwFfhFzJ249Mf

func main() {
	swaps := make(chan swapper.Query)
	swapStatuses := make(chan foundation.SwapStatus)
	swapQueries := make(chan status.Query)

	accounts, err := utils.LoadAccounts("../../secrets/bob.json")
	if err != nil {
		panic(err)
	}

	ldb, err := leveldb.NewStore("../../secrets/bob")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go co.ParBegin(
		func() {
			handler := server.NewHandler(swaps, swapQueries)
			listener, err := net.Listen("tcp", ":18517")
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				if err := http.Serve(listener, handler); err != nil {
					log.Fatal(err)
				}
			}()
			<-done
			listener.Close()
		},
		func() {
			stdLogger := logger.NewStdOut()
			builder := binder.NewBuilder(accounts, stdLogger)
			storage := storage.New(ldb)
			swapper := swapper.New(storage, builder, stdLogger)
			swapper.Run(swaps, swapStatuses, done)
		},
		func() {
			monitor := status.New()
			monitor.Run(swapStatuses, swapQueries, done)
		})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
