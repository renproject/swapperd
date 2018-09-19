package config

import (
	"fmt"

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
func New(net, loc string) config.Config {
	conf := Global
	conf.StoreLocation = loc + "/db"
	conf.RenEx = NewRenExNetwork(net)
	switch net {
	case "nightly":
		conf.RenEx.Orderbook = ""
		conf.RenEx.Settlement = ""
		conf.Ethereum.Swapper = "0xa80c64Cc2c3e29B44CaB2475F6eAd0D523715A4E"
		return conf
	case "falcon":
		conf.RenEx.Orderbook = ""
		conf.RenEx.Settlement = ""
		conf.Ethereum.Swapper = ""
		return conf
	case "testnet":
		conf.RenEx.Orderbook = ""
		conf.RenEx.Settlement = ""
		conf.Ethereum.Swapper = ""
		return conf
	default:
		panic("unimplemented")
	}
}
