package fund

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tyler-smith/go-bip39"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
)

type Config struct {
	Mnemonic string           `json:"mnemonic"`
	Ethereum BlockchainConfig `json:"ethereum"`
	Bitcoin  BlockchainConfig `json:"bitcoin"`
}

type BlockchainConfig struct {
	Network Network `json:"network"`
	Address string  `json:"address"`
	Tokens  []Token `json:"tokens"`
}

type Network struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Token struct {
	Name    string `json:"name"`
	Token   string `json:"erc20"`
	Swapper string `json:"swapper"`
}

func (manager *manager) EthereumAccount(password string) (beth.Account, error) {
	var derivationPath []uint32
	switch manager.config.Ethereum.Network.Name {
	case "kovan", "ropsten":
		derivationPath = []uint32{44, 1, 0, 0, 0}
	case "mainnet":
		derivationPath = []uint32{44, 60, 0, 0, 0}
	}
	privKey, err := manager.loadKey(password, derivationPath)
	if err != nil {
		return nil, err
	}
	ethAccount, err := beth.NewAccount(manager.config.Ethereum.Network.URL, privKey)
	if err != nil {
		return nil, err
	}
	for _, token := range manager.config.Ethereum.Tokens {
		if err := ethAccount.WriteAddress(SwapperKey(token.Name), common.HexToAddress(token.Swapper)); err != nil {
			return nil, err
		}
		if err := ethAccount.WriteAddress(ERC20Key(token.Name), common.HexToAddress(token.Token)); err != nil {
			return nil, err
		}
	}
	return ethAccount, nil
}

// BitcoinAccount returns the bitcoin account
func (manager *manager) BitcoinAccount(password string) (libbtc.Account, error) {
	var derivationPath []uint32
	switch manager.config.Bitcoin.Network.Name {
	case "testnet", "testnet3":
		derivationPath = []uint32{44, 1, 0, 0, 0}
	case "mainnet":
		derivationPath = []uint32{44, 0, 0, 0, 0}
	}
	privKey, err := manager.loadKey(password, derivationPath)
	if err != nil {
		return nil, err
	}
	return libbtc.NewAccount(libbtc.NewBlockchainInfoClient(manager.config.Bitcoin.Network.Name), privKey), nil
}

func (manager *manager) loadKey(password string, path []uint32) (*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(manager.config.Mnemonic, password)
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
