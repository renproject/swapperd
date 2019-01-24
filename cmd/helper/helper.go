package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
)

func main() {
	password := ""
	mnemonic := ""
	derivationPath := []uint32{44, 60, 0, 0, 0}

	seed := bip39.NewSeed(mnemonic, password)
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		panic(err)
	}
	for _, val := range derivationPath {
		key, err = key.NewChildKey(val)
		if err != nil {
			panic(err)
		}
	}
	privKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PrivateKey: %x ", crypto.FromECDSA(privKey))
}
