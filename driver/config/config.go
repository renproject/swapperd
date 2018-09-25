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
}

// RenExMainnet is the RenEx config object on the mainnet
var RenExMainnet = config.RenExNetwork{
	Network:    "mainnet",
	Ingress:    "renex-ingress-mainnet.herokuapp.com",
	Settlement: "0x908262dE0366E42d029B0518D5276762c92B21e1",
	Orderbook:  "0xd5fAEF6b5eE44391FFaa42d732c88b86C73ed287",
}

// RenExTestnet is the RenEx config object on the testnet
var RenExTestnet = config.RenExNetwork{
	Network:    "testnet",
	Ingress:    "renex-ingress-testnet.herokuapp.com",
	Settlement: "0x908262dE0366E42d029B0518D5276762c92B21e1",
	Orderbook:  "0xd5fAEF6b5eE44391FFaa42d732c88b86C73ed287",
}

// RenExNightly is the RenEx config object on the nightly testnet
var RenExNightly = config.RenExNetwork{
	Network:    "nightly",
	Ingress:    "renex-ingress-nightly.herokuapp.com",
	Settlement: "0x5f25233ca99104D31612D4fB937B090d5A2EbB75",
	Orderbook:  "0x376127aDc18260fc238eBFB6626b2F4B59eC9b66",
}

// EthereumMainnet is the ethereum config object on the kovan testnet
var EthereumMainnet = config.EthereumNetwork{
	Network: "mainnet",
	Swapper: "0x6b8bB175c092DE7d81860B18DB360B734A2598e0",
	URL:     "https://mainnet.infura.io",
}

// EthereumKovan is the ethereum config object on the kovan testnet
var EthereumKovan = config.EthereumNetwork{
	Network: "kovan",
	Swapper: "0x9231e9859c8773C17ac896B7fa505AB271F14ea4",
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
		Network: net,
		Ingress: fmt.Sprintf(""),
	}
}

// New creates a new config object from the config data object
func New(loc, net string) (config.Config, error) {
	conf := config.Config{}
	configFilename := fmt.Sprintf("%s/config-%s.json", loc, net)
	data, err := ioutil.ReadFile(configFilename)
	if err == nil {
		if err := json.Unmarshal(data, &conf); err != nil {
			return conf, err
		}
	} else {
		switch net {
		case "nightly", "Nightly":
			conf = NewNightly(loc)
		case "testnet", "Testnet":
			conf = NewTestnet(loc)
		case "", "mainnet", "Mainnet":
			conf = NewMainnet(loc)
		default:
			return conf, fmt.Errorf("Unknown network: %s", net)
		}
	}
	SaveToFile(configFilename, conf)
	return conf, nil
}

func SaveToFile(filename string, conf config.Config) {
	data, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, data, os.FileMode(0644))
}

// NewMainnet creates a new mainnet config
func NewMainnet(loc string) config.Config {
	return config.Config{
		Version:             Global.Version,
		SupportedCurrencies: Global.SupportedCurrencies,
		HomeDir:             loc,
		Ethereum:            EthereumMainnet,
		Bitcoin:             BitcoinMainnet,
		RenEx:               RenExMainnet,
	}
}

// NewTestnet creates a new testnet config
func NewTestnet(loc string) config.Config {
	return config.Config{
		Version:             Global.Version,
		SupportedCurrencies: Global.SupportedCurrencies,
		HomeDir:             loc,
		Ethereum:            EthereumKovan,
		Bitcoin:             BitcoinTestnet,
		RenEx:               RenExTestnet,
	}
}

// NewNightly creates a new nightly config
func NewNightly(loc string) config.Config {
	return config.Config{
		Version:             Global.Version,
		SupportedCurrencies: Global.SupportedCurrencies,
		HomeDir:             loc,
		Ethereum:            EthereumKovan,
		Bitcoin:             BitcoinTestnet,
		RenEx:               RenExNightly,
	}
}
