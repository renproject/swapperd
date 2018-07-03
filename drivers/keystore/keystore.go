package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/republicprotocol/atom-go/adapters/keystore"
)

func main() {
	keyPath := "/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/keystore/keys.json"
	ksPath := "/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/keystore/keystore.json"

	ethNet := flag.String("ethereum", "ganache", "Ethereum Network")
	btcNet := flag.String("bitcoin", "regtest", "Bitcoin Network")

	var keys Keys

	raw, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &keys)
	if err != nil {
		panic(err)
	}

	eth, btc := getKeys(*ethNet, *btcNet, keys)

	ks := keystore.NewKeystore(ksPath)

	err = ks.UpdateEtherumKeyString(eth)
	if err != nil {
		panic(err)
	}
	err = ks.UpdateBitcoinKeyString(btc)
	if err != nil {
		panic(err)
	}
}

type Keys struct {
	Ethereum EthereumKeys `json:"ethereum"`
	Bitcoin  BitcoinKeys  `json:"bitcoin"`
}

type EthereumKeys struct {
	Ganache string `json:"ganache"`
	Ropsten string `json:"ropsten"`
	Kovan   string `json:"kovan"`
	Mainnet string `json:"mainnet"`
}

type BitcoinKeys struct {
	Regtest string `json:"regtest"`
	Testnet string `json:"testnet"`
	Mainnet string `json:"mainnet"`
}

func getKeys(ethereum string, bitcoin string, keys Keys) (string, string) {
	var eth, btc string

	switch ethereum {
	case "ganache":
		eth = keys.Ethereum.Ganache
	case "ropsten":
		eth = keys.Ethereum.Ropsten
	case "kovan":
		eth = keys.Ethereum.Kovan
	case "mainnet":
		eth = keys.Ethereum.Mainnet
	default:
		panic("Unknown Ethereum Network")
	}

	switch bitcoin {
	case "regtest":
		btc = keys.Bitcoin.Regtest
	case "testnet":
		btc = keys.Bitcoin.Testnet
	case "mainnet":
		btc = keys.Bitcoin.Mainnet
	default:
		panic("Unknown Bitcoin Network")
	}

	return eth, btc
}
