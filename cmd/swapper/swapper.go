package main

import (
	"flag"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"runtime"

	guardianAdapter "github.com/republicprotocol/renex-swapper-go/adapter/guardian"
	"github.com/republicprotocol/renex-swapper-go/adapter/http"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	renexAdapter "github.com/republicprotocol/renex-swapper-go/adapter/renex"
	stateAdapter "github.com/republicprotocol/renex-swapper-go/adapter/state"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/driver/config"
	httpDriver "github.com/republicprotocol/renex-swapper-go/driver/http"
	keystoreDriver "github.com/republicprotocol/renex-swapper-go/driver/keystore"
	loggerDriver "github.com/republicprotocol/renex-swapper-go/driver/logger"
	"github.com/republicprotocol/renex-swapper-go/driver/network"
	storeDriver "github.com/republicprotocol/renex-swapper-go/driver/store"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

func main() {
	port := flag.String("port", "18516", "HTTP Atom port")
	repNet := flag.String("network", "mainnet", "Republic Protocol Network")
	keyphrase := flag.String("passphrase", "", "Keyphrase to unlock keystore")
	location := flag.String("loc", getHome()+"/.swapper", "Location of the swapper directory")
	flag.Parse()

	conf, err := config.New(*location, *repNet)
	if err != nil {
		panic(err)
	}

	ks, err := keystoreDriver.LoadFromFile(conf, *keyphrase)
	if err != nil {
		panic(err)
	}

	db, err := storeDriver.NewLevelDB(conf.HomeDir + "/db")
	if err != nil {
		panic(err)
	}

	logger := loggerDriver.NewStdOut()
	state := state.NewState(stateAdapter.New(db, logger))
	ingressNet := network.NewIngress(conf.RenEx.Ingress, ks.GetKey(token.ETH).(keystore.EthereumKey))

	binder, err := renexAdapter.NewBinder(conf, logger)
	if err != nil {
		panic(err)
	}

	renexSwapper := renex.NewRenEx(renexAdapter.New(conf, ks, ingressNet, state, logger, binder))
	guardian := guardian.NewGuardian(guardianAdapter.New(conf, ks, state, logger))

	errCh := make(chan error, 1)
	go renexSwapper.Run(errCh)
	go guardian.Run(errCh)
	go func() {
		for err := range errCh {
			fmt.Println("Swapper Error: ", err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer close(errCh)
		_ = <-c
		log.Println("Stopping the atom box safely")
		os.Exit(1)
	}()

	httpAdapter := http.NewAdapter(conf, ks, renexSwapper)
	log.Println(fmt.Sprintf("0.0.0.0:%s", *port))
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%s", *port), httpDriver.NewServer(httpAdapter)))
}

func getHome() string {
	system := runtime.GOOS
	switch system {
	case "window":
		return os.Getenv("userprofile")
	case "linux", "darwin":
		return os.Getenv("HOME")
	default:
		panic("unknown Operating System")
	}
}
