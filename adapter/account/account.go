package account

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/beth-go"
	libbtc "github.com/republicprotocol/libbtc-go"
	"github.com/tyler-smith/go-bip39"
)

type Config struct {
	Ethereum EthereumConfig `json:"ethereum"`
	Bitcoin  string         `json:"bitcoin"`
}

type EthereumConfig struct {
	URL     string          `json:"url"`
	Network string          `json:"network"`
	Swapper string          `json:"swapper"`
	Tokens  []EthereumToken `json:"tokens"`
}

type EthereumToken struct {
	Name    string `json:"name"`
	ERC20   string `json:"erc20"`
	Swapper string `json:"swapper"`
}

type Accounts interface {
	GetBitcoinAccount(password string) (libbtc.Account, error)
	GetEthereumAccount(password string) (beth.Account, error)
}

type accounts struct {
	mnemonic string
	config   Config
}

type accounts2 struct {
	btcAccount libbtc.Account
	ethAccount beth.Account
}

func New(mnemonic string, config Config) Accounts {
	return &accounts{
		mnemonic: mnemonic,
		config:   config,
	}
}

func (accounts *accounts) GetEthereumAccount(password string) (beth.Account, error) {
	derivationPath := []uint32{}
	switch accounts.config.Ethereum.Network {
	case "kovan", "ropsten":
		derivationPath = []uint32{44, 1, 0, 0}
	case "mainnet":
		derivationPath = []uint32{44, 60, 0, 0}
	}
	privKey, err := accounts.loadKey(password, derivationPath)
	if err != nil {
		return nil, err
	}
	ethAccount, err := beth.NewAccount(accounts.config.Ethereum.URL, privKey)
	if err != nil {
		return nil, err
	}
	ethAccount.WriteAddress(SwapperKey("ETH"), common.HexToAddress(accounts.config.Ethereum.Swapper))
	for _, token := range accounts.config.Ethereum.Tokens {
		if err := ethAccount.WriteAddress(SwapperKey(token.Name), common.HexToAddress(token.Swapper)); err != nil {
			return nil, err
		}
		if err := ethAccount.WriteAddress(ERC20Key(token.Name), common.HexToAddress(token.ERC20)); err != nil {
			return nil, err
		}
	}
	return ethAccount, nil
}

// GetBitcoinAccount returns the bitcoin account
func (accounts *accounts) GetBitcoinAccount(password string) (libbtc.Account, error) {
	derivationPath := []uint32{}
	switch accounts.config.Bitcoin {
	case "testnet", "testnet3":
		derivationPath = []uint32{44, 1, 0, 0, 0}
	case "mainnet":
		derivationPath = []uint32{44, 0, 0, 0, 0}
	}
	privKey, err := accounts.loadKey(password, derivationPath)
	if err != nil {
		return nil, err
	}
	return libbtc.NewAccount(libbtc.NewBlockchainInfoClient(accounts.config.Bitcoin), privKey), nil
}

func (accounts *accounts) loadKey(password string, path []uint32) (*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(accounts.mnemonic, password)
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	for _, val := range path {
		key, err = key.Child(val)
		if err != nil {
			return nil, err
		}
	}
	privKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privKey.ToECDSA(), nil
}

func SwapperKey(token string) string {
	return fmt.Sprintf("SWAPPER:%s", token)
}

func ERC20Key(token string) string {
	return fmt.Sprintf("ERC20:%s", token)
}
