package main

import (
	"flag"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"

	"github.com/republicprotocol/renex-swapper-go/services/logger"
	"github.com/republicprotocol/renex-swapper-go/services/watchdog"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/renex-swapper-go/adapters/atoms"
	"github.com/republicprotocol/renex-swapper-go/adapters/blockchain/binder"
	btcClient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/btc"
	ethClient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/general"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/network"
	"github.com/republicprotocol/renex-swapper-go/adapters/http"
	loggerAdapter "github.com/republicprotocol/renex-swapper-go/adapters/logger"
	"github.com/republicprotocol/renex-swapper-go/adapters/store/leveldb"
	"github.com/republicprotocol/renex-swapper-go/adapters/watchdog/client"
	"github.com/republicprotocol/renex-swapper-go/services/guardian"
	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/republicprotocol/renex-swapper-go/services/watch"
)

type watchAdapter struct {
	atoms.AtomBuilder
	binder.Binder
	watchdog.WatchdogClient
	logger.Logger
}

func main() {
	home := getHome()

	port := flag.String("port", "18516", "HTTP Atom port")
	confPath := flag.String("config", home+"/.swapper/config.json", "Location of the config file")
	keystrPath := flag.String("keystore", home+"/.swapper/keystore.json", "Location of the keystore file")
	networkPath := flag.String("network", home+"/.swapper/network.json", "Location of the network file")

	flag.Parse()

	conf, err := config.LoadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	keystr, err := keystore.Load(*keystrPath)
	if err != nil {
		panic(err)
	}

	log.Println("Swapper is syncing with the bitcoin node, this might take few minutes to complete")
	net, err := network.LoadNetwork(*networkPath)

	dbLoc, err := conf.StoreLocation()
	if err != nil {
		panic(err)
	}

	db, err := leveldb.NewLDBStore(dbLoc)
	if err != nil {
		panic(err)
	}
	state := store.NewState(db, loggerAdapter.NewStdOutLogger())

	watcher, err := buildWatcher(conf, net, keystr, state)
	if err != nil {
		panic(err)
	}

	guardian, err := buildGuardian(net, keystr, state)
	if err != nil {
		panic(err)
	}

	errCh1 := watcher.Start()
	watcher.Notify()

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
		watcher.Stop()
		log.Println("Stopping the guardian service")
		guardian.Stop()
		log.Println("Stopping the atom box safely")
		os.Exit(1)
	}()

	httpAdapter := http.NewBoxHttpAdapter(conf, net, keystr, watcher)
	log.Println(fmt.Sprintf("0.0.0.0:%s", *port))
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%s", *port), http.NewServer(httpAdapter)))

}

func buildGuardian(net network.Config, keystore keystore.Keystore, state store.State) (guardian.Guardian, error) {
	atomBuilder, err := atoms.NewAtomBuilder(net, keystore)
	if err != nil {
		return nil, err
	}
	return guardian.NewGuardian(atomBuilder, state), nil
}

func buildWatcher(gen config.Config, net network.Config, keystore keystore.Keystore, state store.State) (watch.Watch, error) {
	ethConn, err := ethClient.Connect(net)
	if err != nil {
		return nil, err
	}

	btcConn, err := btcClient.Connect(net)
	if err != nil {
		return nil, err
	}

	ethKey, err := keystore.GetKey(1, 0)
	if err != nil {
		return nil, err
	}

	btcKey, err := keystore.GetKey(0, 0)
	if err != nil {
		return nil, err
	}

	_WIF := btcKey.GetKeyString()
	if err != nil {
		return nil, err
	}

	WIF, err := btcutil.DecodeWIF(_WIF)
	if err != nil {
		return nil, err
	}

	err = btcConn.Client.ImportPrivKeyLabel(WIF, "Atom")
	if err != nil {
		return nil, err
	}

	privKey, err := ethKey.GetKey()
	if err != nil {
		return nil, err
	}

	owner := bind.NewKeyedTransactor(privKey)
	owner.GasLimit = 3000000

	ethBinder, err := binder.NewBinder(privKey, ethConn)

	watchdog := client.NewWatchdogHTTPClient(gen)

	atomBuilder, err := atoms.NewAtomBuilder(net, keystore)
	wAdapter := watchAdapter{
		atomBuilder,
		ethBinder,
		watchdog,
		loggerAdapter.NewStdOutLogger(),
	}

	watcher := watch.NewWatch(&wAdapter, state)
	return watcher, nil
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
