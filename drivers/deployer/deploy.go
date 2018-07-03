package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	"github.com/republicprotocol/atom-go/adapters/keystore"
)

func main() {
	ksPath := "/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/keystore/gareg.json"
	// ksPath := "/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/keystore/gareg.json"
	// ksPath := "/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/keystore/gareg.json"
	// ksPath := "/Users/susruth/go/src/github.com/republicprotocol/atom-go/drivers/keystore/gareg.json"

	ks := keystore.NewKeystore(ksPath)

	key, err := ks.LoadKeypair("ethereum")
	if err != nil {
		panic(err)
	}

	auth := bind.NewKeyedTransactor(key)
	err = deployContracts("ganache", auth)
	if err != nil {
		panic(err)
	}
}

func deployContracts(network string, owner *bind.TransactOpts) error {

	ethclient, err := ethclient.Dial(getURL(network))
	if err != nil {
		return err
	}

	// Deploy Atom contract
	AtomAddress, tx, _, err := bindings.DeployAtomSwap(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Network contract
	NetworkAddress, tx, _, err := bindings.DeployAtomNetwork(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Info contract
	InfoAddress, tx, _, err := bindings.DeployAtomInfo(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Wallet contract
	WalletAddress, tx, _, err := bindings.DeployAtomWallet(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	fmt.Println("\"Atom Address\": \"", AtomAddress.Hex(), "\"")
	fmt.Println("\"Network Address\": \"", NetworkAddress.Hex(), "\"")
	fmt.Println("\"Info Address\": \"", InfoAddress.Hex(), "\"")
	fmt.Println("\"Wallet Address\": \"", WalletAddress.Hex(), "\"")
	return nil
}

func getURL(network string) string {
	switch network {
	case "ganache":
		return "http://localhost:8545"
	case "kovan":
		return "https://kovan.infura.io"
	case "ropsten":
		return "https://kovan.infura.io"
	case "mainnet":
		return "https://kovan.infura.io"
	default:
		panic("Unknown Ethereum Network")
	}
}

func deploy(ctx context.Context, network string, client *ethclient.Client, tx *types.Transaction) error {
	switch network {
	case "ganache":
		time.Sleep(100 * time.Millisecond)
		return nil
	default:
		_, err := bind.WaitDeployed(ctx, bind.DeployBackend(client), tx)
		return err
	}
}
