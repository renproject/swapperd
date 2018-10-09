package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/keystore"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/driver/config"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/republicprotocol/swapperd/utils"
)

func buildConfigs() (config.Config, keystore.Keystore, keystore.Keystore) {
	config := configDriver.New("", "nightly")
	keys := utils.LoadTestKeys("../../secrets/test.json")
	btcKeyA, err := keystore.NewBitcoinKey(keys.Alice.Bitcoin, "mainnet")
	if err != nil {
		panic(err)
	}
	ethPrivKeyA, err := crypto.HexToECDSA(keys.Alice.Ethereum)
	if err != nil {
		panic(err)
	}
	ethKeyA, err := keystore.NewEthereumKey(ethPrivKeyA, "kovan")
	if err != nil {
		panic(err)
	}
	btcKeyB, err := keystore.NewBitcoinKey(keys.Bob.Bitcoin, "mainnet")
	if err != nil {
		panic(err)
	}
	ethPrivKeyB, err := crypto.HexToECDSA(keys.Bob.Ethereum)
	if err != nil {
		panic(err)
	}
	ethKeyB, err := keystore.NewEthereumKey(ethPrivKeyB, "kovan")
	if err != nil {
		panic(err)
	}
	ksA := keystore.New(btcKeyA, ethKeyA)
	ksB := keystore.New(btcKeyB, ethKeyB)
	return config, ksB, ksA
}

func buildRequests(ksA, ksB keystore.Keystore) (foundation.Swap, foundation.Swap) {

	var aliceSwapID, bobSwapID foundation.SwapID
	var aliceSecret [32]byte
	rand.Read(aliceSwapID[:])
	rand.Read(bobSwapID[:])
	rand.Read(aliceSecret[:])
	aliceSecretHash := sha256.Sum256(aliceSecret[:])
	timelock := time.Now().Unix() + 48*60*60

	aliceReq := foundation.Swap{
		ID:                 aliceSwapID,
		TimeLock:           timelock,
		Secret:             aliceSecret,
		SecretHash:         aliceSecretHash,
		SendToAddress:      ksB.GetKey(foundation.TokenWBTC).(keystore.EthereumKey).Address.String(),
		ReceiveFromAddress: ksB.GetKey(foundation.TokenBTC).(keystore.BitcoinKey).AddressString,
		SendValue:          big.NewInt(20000),
		ReceiveValue:       big.NewInt(20000),
		SendToken:          foundation.TokenWBTC,
		ReceiveToken:       foundation.TokenBTC,
		IsFirst:            true,
	}

	areq, err := json.Marshal(aliceReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("Alice Request: ", hex.EncodeToString(areq))

	bobReq := foundation.Swap{
		UID:                bobSwapID,
		TimeLock:           timelock,
		Secret:             [32]byte{},
		SecretHash:         aliceSecretHash,
		SendToAddress:      ksA.GetKey(foundation.TokenBTC).(keystore.BitcoinKey).AddressString,
		ReceiveFromAddress: ksA.GetKey(foundation.TokenWBTC).(keystore.EthereumKey).Address.String(),
		SendValue:          big.NewInt(20000),
		ReceiveValue:       big.NewInt(20000),
		SendToken:          foundation.TokenBTC,
		ReceiveToken:       foundation.TokenWBTC,
		IsFirst:            false,
	}

	breq, err := json.Marshal(bobReq)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println("Bob Request: ", hex.EncodeToString(breq))

	return aliceReq, bobReq
}

func buildSwaps() {
	conf, ksA, ksB := buildConfigs()
	reqA, reqB := buildRequests(ksA, ksB)
	aliceSwapper, bobSwapper := buildSwappers(conf, ksA, ksB)
	aliceSwap, err := aliceSwapper.NewSwap(reqA)
	Expect(err).Should(BeNil())
	bobSwap, err := bobSwapper.NewSwap(reqB)
	Expect(err).Should(BeNil())
	return aliceSwap, bobSwap
}

func main() {
	conf, aliceKS, bobKS := buildConfigs()
	aliceSwapper, bobSwapper := buildSwaps()
	co.ParBegin(
		func() {
			swapper.Swap()
		},
		func() {
			swapper.Swap()
		},
	)
}
