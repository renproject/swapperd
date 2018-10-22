package config

import (
	"github.com/republicprotocol/swapperd/adapter/account"
)

// Testnet is the testnet config object
var Testnet = account.Config{
	Bitcoin: "testnet",
	Ethereum: account.EthereumConfig{
		URL:     "https://kovan.infura.io",
		Network: "kovan",
		Swapper: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
		Tokens: []account.EthereumToken{
			account.EthereumToken{
				Name:    "WBTC",
				ERC20:   "0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0",
				Swapper: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
			},
		},
	},
}
