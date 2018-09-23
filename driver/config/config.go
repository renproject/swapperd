package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
)

// Global is the default global config
var Global = config.Config{
	Version:             "0.1.0",
	SupportedCurrencies: []string{"ETH", "BTC"},
	Ethereum:            EthereumKovan,
	Bitcoin:             BitcoinMainnet,
}

// EthereumKovan is the ethereum config object on the kovan testnet
var EthereumKovan = config.EthereumNetwork{
	Network: "kovan",
	URL:     "https://kovan.infura.io",
}

// BitcoinTestnet is the bitcoin config object on the bitcoin testnet
var BitcoinTestnet = config.BitcoinNetwork{
	Network: "testnet",
	URL:     "https://testnet.blockchain.info",
}

// BitcoinMainnet is the bitcoin config object on the bitcoin mainnet
var BitcoinMainnet = config.BitcoinNetwork{
	Network: "mainnet",
	URL:     "https://blockchain.info",
}

// NewRenExNetwork creates a RenEx config object for the given RenEx network
func NewRenExNetwork(net string) config.RenExNetwork {
	return config.RenExNetwork{
		Network:  net,
		Watchdog: fmt.Sprintf("renex-watchdog-%s.herokuapp.com", net),
		Ingress:  fmt.Sprintf("renex-ingress-nightly.herokuapp.com"),
	}
}

// New creates a new config object from the config data object
func New(loc, net string) config.Config {
	conf := config.Config{}
	filename := fmt.Sprintf("%s/config-%s.json", loc, net)
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(data, &conf); err == nil {
			return conf
		}
	}
	conf = Global
	conf.StoreLocation = loc + "/db"
	conf.RenEx = NewRenExNetwork(net)
	switch net {
	case "nightly":
		conf.RenEx.Orderbook = ""
		conf.RenEx.Settlement = ""
		conf.Ethereum.Swapper = "0x9231e9859c8773C17ac896B7fa505AB271F14ea4"
	case "falcon":
		conf.RenEx.Orderbook = ""
		conf.RenEx.Settlement = ""
		conf.Ethereum.Swapper = ""
	case "testnet":
		conf.RenEx.Orderbook = ""
		conf.RenEx.Settlement = ""
		conf.Ethereum.Swapper = ""
	default:
		panic("Unknown republic network" + net)
	}
	SaveToFile(filename, conf)
	return conf
}

func SaveToFile(filename string, conf config.Config) {
	data, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, data, os.FileMode(0644))
}
