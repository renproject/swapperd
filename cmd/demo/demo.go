package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/republicprotocol/libbtc-go"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/btc"
	configAdapter "github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/eth/client"
	"github.com/republicprotocol/swapperd/adapter/eth/erc20"
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

	// 20000
	value := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 32}

	// 0.5 BTC
	// value := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 250, 240, 128}

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

	return aliceReq, bobReq
}

func buildBinders(config configAdapter.Config, ks keystore.Keystore, logger core.Logger, req foundation.Swap) (core.SwapContractBinder, core.SwapContractBinder) {
	ethereumClient, err := client.New(client.NetworkConfig{
		URL:     config.Ethereum.URL,
		Network: config.Ethereum.Network,
		Tokens: []client.EthereumToken{
			client.EthereumToken{
				Name:           "WBTC",
				TokenAddress:   "0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0",
				SwapperAddress: "0x2218Fa20c33765e7e01671eE6AacA75FbAf3A974",
			},
		},
	}, ks.GetKey(foundation.TokenWBTC).(keystore.EthereumKey).PrivateKey)
	if err != nil {
		panic(err)
	}
	wbtcBinder, err := erc20.NewERC20Atom(ethereumClient, logger, req)
	if err != nil {
		panic(err)
	}

	btcAccount := libbtc.NewAccount(libbtc.NewBlockchainInfoClient("testnet"), ks.GetKey(foundation.TokenBTC).(keystore.BitcoinKey).PrivateKey.ToECDSA())
	btcBinder, err := btc.NewBitcoinAtom(btcAccount, logger, req)
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
					if !res.Success {
						fmt.Println("Atomic Swap Failed!!!")
						continue
					}
					logger.LogDebug(res.ID, "Atomic Swap Successful")
				}
			}
		},
	)
}
