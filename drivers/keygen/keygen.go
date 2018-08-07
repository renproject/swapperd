package main

import (
	"encoding/hex"
	"flag"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/keystore"
)

func main() {
	ethNet := flag.String("ethereum", "kovan", "Which ethereum network to use")
	btcNet := flag.String("bitcoin", "testnet", "Which bitcoin network to use")

	keyPair, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	keyStr := hex.EncodeToString(crypto.FromECDSA(keyPair))
	BTCKey, err := keystore.NewBitcoinKey(keyStr, *btcNet)
	ETHKey, err := keystore.NewEthereumKey(keyStr, *ethNet)
	f, err := os.Create(os.Getenv("HOME") + "/.swapper/keystore.json")
	f.Close()
	ks, err := keystore.LoadKeystore(os.Getenv("HOME") + "/.swapper/keystore.json")
	if err != nil {
		panic(err)
	}
	ks.EthereumKey = ETHKey
	ks.BitcoinKey = BTCKey
	if err := ks.Update(); err != nil {
		panic(err)
	}
}
