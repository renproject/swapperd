package main

import (
	"github.com/republicprotocol/atom-go/adapters/config"
)

func main() {
}

func getEthConfig(network string) config.EthereumConfig {
	switch network {
	case "ganache":
		return config.EthereumConfig
	}
}
