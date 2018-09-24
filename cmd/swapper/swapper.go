package main

import (
	"flag"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"

	guardianAdapter "github.com/republicprotocol/renex-swapper-go/adapter/guardian"
	"github.com/republicprotocol/renex-swapper-go/adapter/http"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	renexAdapter "github.com/republicprotocol/renex-swapper-go/adapter/renex"
	stateAdapter "github.com/republicprotocol/renex-swapper-go/adapter/state"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	configDriver "github.com/republicprotocol/renex-swapper-go/driver/config"
	httpDriver "github.com/republicprotocol/renex-swapper-go/driver/http"
	keystoreDriver "github.com/republicprotocol/renex-swapper-go/driver/keystore"
	loggerDriver "github.com/republicprotocol/renex-swapper-go/driver/logger"
	"github.com/republicprotocol/renex-swapper-go/driver/network"
	storeDriver "github.com/republicprotocol/renex-swapper-go/driver/store"
	watchdogDriver "github.com/republicprotocol/renex-swapper-go/driver/watchdog"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

func main() {
	port := flag.String("port", "18516", "HTTP Atom port")
	repNet := flag.String("network", "testnet", "Republic Protocol Network")
	keyphrase := flag.String("passphrase", "", "Keyphrase to unlock keystore")
	location := flag.String("loc", getHome()+"/.swapper", "Location of the swapper directory")
	flag.Parse()

	conf := configDriver.New(*location, *repNet)
	ks, err := keystoreDriver.LoadFromFile(*repNet, *location, *keyphrase)
	if err != nil {
		panic(err)
	}
	db, err := storeDriver.NewLevelDB(conf.StoreLocation)
	if err != nil {
		panic(err)
	}

	fmt.Println(ks.GetKey(token.ETH).(keystore.EthereumKey).Address.String())
	fmt.Println(ks.GetKey(token.BTC).(keystore.BitcoinKey).AddressString)

	logger := loggerDriver.NewStdOut()
	state := state.NewState(stateAdapter.New(db, logger))
	ingressNet := network.NewIngress(conf.RenEx.Ingress, ks.GetKey(token.ETH).(keystore.EthereumKey))
	nopWatchdog := watchdogDriver.NewMock()

	binder, err := renexAdapter.NewBinder(conf)
	if err != nil {
		panic(err)
	}

	renexSwapper := renex.NewRenEx(renexAdapter.New(conf, ks, ingressNet, nopWatchdog, state, logger, binder))
	guardian := guardian.NewGuardian(guardianAdapter.New(conf, ks, state, logger))

	errCh1 := renexSwapper.Start()
	renexSwapper.Notify()

	errCh2 := guardian.Start()
	guardian.Notify()

	go func() {
		for err := range errCh1 {
			log.Println("Watcher Error :", err)
		}
	}()

	go func() {
		for err := range errCh2 {
			log.Println("Guardian Error :", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Stopping the swapper service")
		renexSwapper.Stop()
		log.Println("Stopping the guardian service")
		guardian.Stop()
		log.Println("Stopping the atom box safely")
		os.Exit(1)
	}()

	httpAdapter := http.NewAdapter(conf, ks, renexSwapper)
	log.Println(fmt.Sprintf("0.0.0.0:%s", *port))
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%s", *port), httpDriver.NewServer(httpAdapter)))
}

func getHome() string {
	winHome := os.Getenv("userprofile")
	unixHome := os.Getenv("HOME")

	if winHome != "" {
		return winHome
	}

	if unixHome != "" {
		return unixHome
	}

	panic("unknown Operating System")
}
