package main

import (
	"github.com/republicprotocol/atom-go/adapters/configs/keystore"
)

func main() {

	keystore.NewKeystore([]uint32{0, 1}, []string{"regtest", "ganache"}, "./../../secrets/local.alice.json")
	keystore.NewKeystore([]uint32{0, 1}, []string{"regtest", "ganache"}, "./../../secrets/local.bob.json")
	keystore.NewKeystore([]uint32{0, 1}, []string{"testnet", "kovan"}, "./../../secrets/test.alice.json")
	keystore.NewKeystore([]uint32{0, 1}, []string{"testnet", "kovan"}, "./../../secrets/test.bob.json")
	keystore.NewKeystore([]uint32{0, 1}, []string{"mainnet", "mainnet"}, "./../../secrets/main.alice.json")
	keystore.NewKeystore([]uint32{0, 1}, []string{"mainnet", "mainnet"}, "./../../secrets/main.bob.json")

	// keyPair, err := crypto.GenerateKey()
	// if err != nil {
	// 	panic(err)
	// }
	// keyStr := hex.EncodeToString(crypto.FromECDSA(keyPair))
	// BTCKey, err := keystore.NewBitcoinKey(keyStr, *btcNet)
	// ETHKey, err := keystore.NewEthereumKey(keyStr, *ethNet)
	// f, err := os.Create(os.Getenv("HOME") + "/.swapper/keystore.json")
	// f.Close()
	// ks, err := keystore.LoadKeystore(os.Getenv("HOME") + "/.swapper/keystore.json")
	// if err != nil {
	// 	panic(err)
	// }
	// ks.EthereumKey = ETHKey
	// ks.BitcoinKey = BTCKey
	// if err := ks.Update(); err != nil {
	// 	panic(err)
	// }
}
