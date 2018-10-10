package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/btc"
	configAdapter "github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/erc20"
	"github.com/republicprotocol/swapperd/adapter/keystore"
	"github.com/republicprotocol/swapperd/core"
	"github.com/republicprotocol/swapperd/driver/config"
	"github.com/republicprotocol/swapperd/driver/logger"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/republicprotocol/swapperd/utils"
)

func buildConfigs() (configAdapter.Config, keystore.Keystore, keystore.Keystore) {
	config, err := config.New("", "nightly")
	if err != nil {
		panic(err)
	}
	keys := utils.LoadTestKeys("../../secrets/test.json")
	btcKeyA, err := keystore.NewBitcoinKey(keys.Alice.Bitcoin, "testnet")
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
	btcKeyB, err := keystore.NewBitcoinKey(keys.Bob.Bitcoin, "testnet")
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

	value := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 32}

	aliceReq := foundation.Swap{
		ID:                 aliceSwapID,
		TimeLock:           timelock,
		Secret:             aliceSecret,
		SecretHash:         aliceSecretHash,
		SendToAddress:      ksB.GetKey(foundation.TokenWBTC).(keystore.EthereumKey).Address.String(),
		ReceiveFromAddress: ksB.GetKey(foundation.TokenBTC).(keystore.BitcoinKey).AddressString,
		SendValue:          value,
		ReceiveValue:       value,
		SendToken:          foundation.TokenWBTC,
		ReceiveToken:       foundation.TokenBTC,
		IsFirst:            true,
	}

	fmt.Println("Alice Request: ", aliceReq)

	bobReq := foundation.Swap{
		ID:                 bobSwapID,
		TimeLock:           timelock,
		Secret:             [32]byte{},
		SecretHash:         aliceSecretHash,
		SendToAddress:      ksA.GetKey(foundation.TokenBTC).(keystore.BitcoinKey).AddressString,
		ReceiveFromAddress: ksA.GetKey(foundation.TokenWBTC).(keystore.EthereumKey).Address.String(),
		SendValue:          value,
		ReceiveValue:       value,
		SendToken:          foundation.TokenBTC,
		ReceiveToken:       foundation.TokenWBTC,
		IsFirst:            false,
	}
	fmt.Println("Bob Request: ", bobReq)

	return aliceReq, bobReq
}

func buildBinders(config configAdapter.Config, ks keystore.Keystore, logger core.Logger, req foundation.Swap) (core.SwapContractBinder, core.SwapContractBinder) {
	wbtcBinder, err := erc20.NewERC20Atom(config.Ethereum, ks.GetKey(foundation.TokenWBTC).(keystore.EthereumKey), logger, req)
	if err != nil {
		panic(err)
	}

	btcBinder, err := btc.NewBitcoinAtom(config.Bitcoin, ks.GetKey(foundation.TokenBTC).(keystore.BitcoinKey), logger, req)
	if err != nil {
		panic(err)
	}

	if req.SendToken == foundation.TokenBTC {
		return btcBinder, wbtcBinder
	}
	return wbtcBinder, btcBinder
}

func main() {
	conf, aliceKS, bobKS := buildConfigs()
	aliceReq, bobReq := buildRequests(aliceKS, bobKS)
	logger := logger.NewStdOut()
	aliceNativeBinder, aliceForeignBinder := buildBinders(conf, aliceKS, logger, aliceReq)
	bobNativeBinder, bobForeignBinder := buildBinders(conf, bobKS, logger, bobReq)

	results := make(chan core.Result, 2)

	co.ParBegin(
		func() {
			core.Swap(aliceNativeBinder, aliceForeignBinder, logger, aliceReq, results)
		},
		func() {
			core.Swap(bobNativeBinder, bobForeignBinder, logger, bobReq, results)
		},
		func() {
			for i := 0; i < 2; i++ {
				select {
				case res := <-results:
					fmt.Println(res)
				}
			}
		},
	)
}
