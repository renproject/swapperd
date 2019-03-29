package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/renproject/libbtc-go"
	"github.com/renproject/libeth-go"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func (wallet *wallet) EthereumAccount(password string) (libeth.Account, error) {
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
	ethClient, err := wallet.ethereumClient()
	if err != nil {
		return nil, err
	}
	ethAccount, err := libeth.NewAccount(ethClient, privKey)
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
	logger := wallet.logger.WithField("token", "bitcoin")
	logger = logger.WithField("network", wallet.config.Bitcoin.Network.Name)
	client, err := wallet.bitcoinClient()
	if err != nil {
		return nil, err
	}
	return libbtc.NewAccount(client, privKey, logger), nil
}

func (wallet *wallet) ethereumClient() (libeth.Client, error) {
	return libeth.NewInfuraClient(wallet.config.Ethereum.Network.Name, "172978c53e244bd78388e6d50a4ae2fa")
}

func (wallet *wallet) bitcoinClient() (libbtc.Client, error) {
	switch wallet.config.Bitcoin.Network.Name {
	case "mainnet":
		return libbtc.NewBlockchainInfoClient("mainnet")
	case "testnet", "testnet3":
		return libbtc.NewMercuryClient("testnet")
	default:
		return nil, fmt.Errorf("unsupported network: %s", wallet.config.Bitcoin.Network.Name)
	}
}

func (wallet *wallet) loadECDSAKey(password string, path []uint32) (*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(wallet.config.Mnemonic, password)
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	for _, val := range path {
		key, err = key.NewChildKey(val)
		if err != nil {
			return nil, err
		}
	}
	return crypto.ToECDSA(key.Key)
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
