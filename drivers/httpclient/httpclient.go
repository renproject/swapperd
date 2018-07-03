package main

import (
	"crypto/rand"
	"log"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/republicprotocol/atom-go/adapters/http"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/utils"
)

func main() {
	keystr := keystore.NewKeystore("/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/httpclient/keystore.json")

	keys, err := keystr.LoadECDSA()

	if err != nil {
		panic(err)
	}

	key := keys[0]
	orderID := [32]byte{}

	_, err = rand.Read(orderID[:])
	if err != nil {
		panic(err)
	}

	sig, err := ethCrypto.Sign(orderID[:], key)
	if err != nil {
		panic(err)
	}

	sig65, err := utils.ToBytes65(sig)
	if err != nil {
		panic(err)
	}

	log.Println("address", ethCrypto.PubkeyToAddress(key.PublicKey).String())
	log.Println("order id", http.MarshalOrderID(orderID))
	log.Println("sig 65", http.MarshalSignature(sig65))
}
