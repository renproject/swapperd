package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/tyler-smith/go-bip39"
)

func (wallet *wallet) EthereumAccount(password string) (beth.Account, error) {
	var derivationPath []uint32
	switch wallet.config.Ethereum.Network.Name {
	case "kovan", "ropsten":
		derivationPath = []uint32{44, 1, 0, 0, 0}
	case "mainnet":
		derivationPath = []uint32{44, 60, 0, 0, 0}
	}
	privKey, err := wallet.loadECDSAKey(password, derivationPath)
	if err != nil {
		return nil, err
	}
	ethAccount, err := beth.NewAccount(wallet.config.Ethereum.Network.URL, privKey)
	if err != nil {
		return nil, err
	}
	return ethAccount, nil
}

// BitcoinAccount returns the bitcoin account
func (wallet *wallet) BitcoinAccount(password string) (libbtc.Account, error) {
	var derivationPath []uint32
	switch wallet.config.Bitcoin.Network.Name {
	case "testnet", "testnet3":
		derivationPath = []uint32{44, 1, 0, 0, 0}
	case "mainnet":
		derivationPath = []uint32{44, 0, 0, 0, 0}
	}
	privKey, err := wallet.loadECDSAKey(password, derivationPath)
	if err != nil {
		return nil, err
	}
	return libbtc.NewAccount(libbtc.NewBlockchainInfoClient(wallet.config.Bitcoin.Network.Name), privKey), nil
}

func (wallet *wallet) loadECDSAKey(password string, path []uint32) (*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(wallet.config.Mnemonic, password)
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

func (wallet *wallet) ECDSASigner(password string) (ECDSASigner, error) {
	privKey, err := wallet.loadECDSAKey(password, []uint32{0})
	if err != nil {
		return nil, err
	}
	return &ecdsaSigner{privKey}, nil
}

func (wallet *wallet) loadRSAKey(password string) (*rsa.PrivateKey, error) {
	seed := bip39.NewSeed(wallet.config.Mnemonic, password)
	return rsa.GenerateKey(bytes.NewReader(seed), 2048)
}
