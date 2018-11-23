package fund

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
)

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
