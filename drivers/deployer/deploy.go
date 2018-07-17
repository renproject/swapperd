package main

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/owner"
)

func main() {

	var aPath = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/local/configA.json"
	var bPath = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/local/configB.json"
	var ownPath = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/owner.json"

	aTest, err := config.LoadConfig(aPath)
	if err != nil {
		panic(err)
	}

	bTest, err := config.LoadConfig(bPath)
	if err != nil {
		panic(err)
	}

	own, err := owner.LoadOwner(ownPath)
	if err != nil {
		panic(err)
	}

	key, err := crypto.HexToECDSA(own.Ganache)
	if err != nil {
		panic(err)
	}

	auth := bind.NewKeyedTransactor(key)
	err = deployContracts(aTest, auth)
	if err != nil {
		panic(err)
	}
	bTest.SetEthereumConfig(aTest.GetEthereumConfig())
}

func deployContracts(config config.Config, owner *bind.TransactOpts) error {

	network := config.Ethereum.Chain

	ethclient, err := ethclient.Dial(config.Ethereum.URL)
	if err != nil {
		return err
	}

	// Deploy Atom contract
	AtomAddress, tx, _, err := bindings.DeployAtomicSwap(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Info contract
	InfoAddress, tx, _, err := bindings.DeployAtomicInfo(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Info contract
	RENAddress, tx, _, err := bindings.DeployRepublicToken(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Dark Node Registry contract
	DNRAddress, tx, _, err := bindings.DeployDarknodeRegistry(owner, ethclient, RENAddress, big.NewInt(0), big.NewInt(8), big.NewInt(0))
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Order Book contract
	OBAddress, tx, _, err := bindings.DeployOrderbook(owner, ethclient, big.NewInt(0), RENAddress, DNRAddress)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy RenEx Tokens contract
	RenExTokensAddress, tx, _, err := bindings.DeployRenExTokens(owner, ethclient)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Reward Vault contract
	RewardVaultAddress, tx, _, err := bindings.DeployRewardVault(owner, ethclient, DNRAddress)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy RenEx Balances contract
	RenExBalancesAddress, tx, _, err := bindings.DeployRenExBalances(owner, ethclient, RewardVaultAddress)
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	// Deploy Wallet contract
	WalletAddress, tx, _, err := bindings.DeployRenExSettlement(owner, ethclient, OBAddress, RenExTokensAddress, RenExBalancesAddress, big.NewInt(100))
	if err != nil {
		return err
	}

	if err := deploy(context.Background(), network, ethclient, tx); err != nil {
		return err
	}

	ethCfg := config.GetEthereumConfig()
	ethCfg.AtomAddress = AtomAddress.Hex()
	ethCfg.InfoAddress = InfoAddress.Hex()
	ethCfg.WalletAddress = WalletAddress.Hex()
	ethCfg.RepublicTokenAddress = WalletAddress.Hex()
	ethCfg.DarkNodeRegistryAddress = DNRAddress.Hex()
	ethCfg.OrderBookAddress = OBAddress.Hex()
	ethCfg.RenExTokens = RenExTokensAddress.Hex()
	ethCfg.RewardVault = RewardVaultAddress.Hex()
	ethCfg.RenExBalances = RenExBalancesAddress.Hex()
	config.SetEthereumConfig(ethCfg)
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
