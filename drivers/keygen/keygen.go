package main

import (
	"encoding/hex"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/atom-go/adapters/configs/keystore"
)

func main() {
	keyPair, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	keyStr := hex.EncodeToString(crypto.FromECDSA(keyPair))
	BTCKey, err := keystore.NewBitcoinKey(keyStr, "regtest")
	ETHKey, err := keystore.NewEthereumKey(keyStr, "ganache")
	ks, err := keystore.LoadKeystore(os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/local/keystoreB.json")
	if err != nil {
		panic(err)
	}
	ks.EthereumKey = ETHKey
	ks.BitcoinKey = BTCKey
	if err := ks.Update(); err != nil {
		panic(err)
	}
}
