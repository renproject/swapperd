package main

import (
	"flag"
	"fmt"
	"log"
	netHttp "net/http"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/atom-go/adapters/atoms"
	btcClient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethClient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/http"
	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"
	"github.com/republicprotocol/atom-go/adapters/store/leveldb"
	wal "github.com/republicprotocol/atom-go/adapters/wallet/eth"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/atom-go/services/watch"
)

func main() {

	port := flag.String("port", "18516", "HTTP Atom port")
	confPath := flag.String("config", "../config.json", "Location of the config file")
	keystrPath := flag.String("keystore", "../keystore.json", "Location of the keystore file")

	flag.Parse()

	conf, err := config.LoadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	keystr := keystore.NewKeystore(*keystrPath)

	watcher, err := buildWatcher(conf, keystr)
	if err != nil {
		panic(err)
	}

	httpAdapter := http.NewBoxHttpAdapter(conf, keystr, watcher)
	log.Println(fmt.Sprintf("0.0.0.0:%s", *port))
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%s", *port), http.NewServer(httpAdapter)))
}

func buildWatcher(config config.Config, kstr swap.Keystore) (watch.Watch, error) {
	ethConn, err := ethClient.Connect(config)
	if err != nil {
		return nil, err
	}

	btcConn, err := btcClient.Connect(config)
	if err != nil {
		return nil, err
	}

	keys, err := kstr.LoadKeys()
	if err != nil {
		return nil, err
	}

	ethKey := keys[0]
	btcKey := keys[1]

	_WIF, err := btcKey.GetKeyString()
	if err != nil {
		return nil, err
	}

	WIF, err := btcutil.DecodeWIF(_WIF)
	if err != nil {
		return nil, err
	}

	err = btcConn.Client.ImportPrivKey(WIF)
	if err != nil {
		return nil, err
	}

	owner := bind.NewKeyedTransactor(ethKey.GetKey())
	owner.GasLimit = 3000000

	ethNet, err := net.NewEthereumNetwork(ethConn, owner)
	if err != nil {
		return nil, err
	}

	ethInfo, err := ax.NewEthereumAtomInfo(ethConn, owner)
	if err != nil {
		return nil, err
	}

	ethWallet, err := wal.NewEthereumWallet(ethConn, *owner)
	if err != nil {
		return nil, err
	}

	atomBuilder := atoms.NewAtomBuilder(config, keys)

	db, err := leveldb.NewLDBStore(config.StoreLocation())
	if err != nil {
		return nil, err
	}
	str := store.NewSwapState(db)
	watcher := watch.NewWatch(ethNet, ethInfo, ethWallet, atomBuilder, str)
	return watcher, nil
}
