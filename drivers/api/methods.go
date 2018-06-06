package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/atom-go/adapters/eth"
)

func InitiateAtomicSwap(swap AtomicSwap) {

	swapID := make([]byte, 32)
	rand.Read(swapID)
	swapID32, err := toByte32(swapID)
	if err != nil {

	}

	secret := make([]byte, 32)
	rand.Read(swapID)
	secret32, err := toByte32(secret)
	if err != nil {

	}

	value := new(big.Int)
	value, ok := value.SetString(swap.Value, 10)
	if !ok {
		// errors.new("Failed to parse value")
	}

	myOrderID := base58.Decode(swap.MyOrderID)
	matchingOrderID := base58.Decode(swap.MatchingOrderID)
	keyPair, err := crypto.HexToECDSA(swap.PrivateKey)
	auth := bind.NewKeyedTransactor(keyPair)
	conn, err := eth.Connect(eth.Network(swap.Network))

	if err != nil {

	}

	if goFirst(myOrderID, matchingOrderID) {
		atom, err := eth.NewEthereumAtom(context.Background(), conn, auth, swapID32)
		if err != nil {

		}
		atom.Initiate(secret32, auth.From.Bytes(), auth.From.Bytes(), value, time.Now().Add(48*time.Hour).Unix())
	}
}
