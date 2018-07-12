package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcutil"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/adapters/owner"
	wal "github.com/republicprotocol/atom-go/adapters/wallet/eth"
	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/swap"
)

func main() {

	var aliceOrderID, bobOrderID [32]byte

	var conf = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/test/configA.json"
	var keyA = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/test/keystoreA.json"
	var keyB = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/test/keystoreB.json"
	var ownPath = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/owner.json"

	ksA := keystore.NewKeystore(keyA)
	ksB := keystore.NewKeystore(keyB)

	own, err := owner.LoadOwner(ownPath)
	if err != nil {
		panic(err)
	}

	ownerECDSA, err := crypto.HexToECDSA(own.Kovan)
	if err != nil {
		panic(err)
	}
	owner := bind.NewKeyedTransactor(ownerECDSA)

	config, err := config.LoadConfig(conf)
	if err != nil {
		panic(err)
	}

	ethConn, err := ethclient.Connect(config)
	if err != nil {
		panic(err)
	}

	btcConn, err := btcclient.Connect(config)
	if err != nil {
		panic(err)
	}

	err = deposit(ksA, ksB, owner, ethConn, btcConn)
	if err != nil {
		panic(err)
	}

	rand.Read(aliceOrderID[:])
	rand.Read(bobOrderID[:])

	aliceSendValue := big.NewInt(10000000)
	bobSendValue := big.NewInt(10000000)

	aliceCurrency := uint32(1)
	bobCurrency := uint32(0)

	atomMatch := match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, bobSendValue, aliceCurrency, bobCurrency)
	mockWallet, err := wal.NewEthereumWallet(ethConn, *owner)
	if err != nil {
		panic(err)
	}

	err = mockWallet.SetMatch(atomMatch)
	if err != nil {
		panic(err)
	}

	messageAlice := append([]byte("Republic Protocol: open: "), aliceOrderID[:]...)
	signatureDataAlice := ethCrypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(messageAlice))), messageAlice)

	sigAlice, err := ethCrypto.Sign(signatureDataAlice[:], ownerECDSA)
	if err != nil {
		panic(err)
	}

	messageBob := append([]byte("Republic Protocol: open: "), bobOrderID[:]...)
	signatureDataBob := ethCrypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(messageBob))), messageBob)

	sigBob, err := ethCrypto.Sign(signatureDataBob[:], ownerECDSA)
	if err != nil {
		panic(err)
	}

	fmt.Println("Signer:", owner.From.String())
	fmt.Println("Alice:", hex.EncodeToString(aliceOrderID[:]), "Signature:", hex.EncodeToString(sigAlice))
	fmt.Println("Bob:", hex.EncodeToString(bobOrderID[:]), "Signature:", hex.EncodeToString(sigBob))
}

func deposit(ksA, ksB swap.Keystore, auth *bind.TransactOpts, ethConn ethclient.Conn, btcConn btcclient.Conn) error {
	Alice, err := ksA.LoadKeys()
	if err != nil {
		return err
	}
	Bob, err := ksB.LoadKeys()
	if err != nil {
		return err
	}

	aEth := Alice[0]
	aBtc := Alice[1]
	bEth := Bob[0]
	bBtc := Bob[1]

	aliceEth, err := aEth.GetAddress()
	if err != nil {
		return err
	}

	aliceBtc, err := aBtc.GetAddress()
	if err != nil {
		return err
	}

	bobEth, err := bEth.GetAddress()
	if err != nil {
		return err
	}

	bobBtc, err := bBtc.GetAddress()
	if err != nil {
		return err
	}

	aliceEthAddr := common.BytesToAddress(aliceEth)
	bobEthAddr := common.BytesToAddress(bobEth)
	// ethConn.Transfer(aliceAddr, auth, 1000000000000)

	bobBtcAddr, err := btcutil.DecodeAddress(string(bobBtc), btcConn.ChainParams)
	if err != nil {
		return err
	}

	aliceBtcAddr, err := btcutil.DecodeAddress(string(aliceBtc), btcConn.ChainParams)
	if err != nil {
		return err
	}

	// bobVal, err := btcutil.NewAmount(0.05)
	// if err != nil {
	// 	return err
	// }

	// _, err = btcConn.Client.SendToAddress(bobAddr, bobVal)

	fmt.Println("Deposit Ether to Alice:", aliceEthAddr.String())
	fmt.Println("Deposit Ether to Bob:", bobEthAddr.String())
	fmt.Println("Deposit Bitcoin to Alice:", aliceBtcAddr)
	fmt.Println("Deposit Bitcoin to Bob:", bobBtcAddr)

	return nil
}
