package keystore

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type bitcoinKey key

func (key *bitcoinKey) GetAddress() ([]byte, error) {
	chainParams, err := getChainParams(key.Network)
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.DecodeWIF(key.PrivateKey)
	if err != nil {
		return nil, err
	}

	serializedPubKey := wif.SerializePubKey()

	pubKey, err := btcutil.NewAddressPubKey(serializedPubKey, chainParams)
	if err != nil {
		return nil, err
	}

	return []byte(pubKey.EncodeAddress()), nil
}

func (key *bitcoinKey) GetKeyString() string {
	return key.PrivateKey
}

func (key *bitcoinKey) GetKey() (*ecdsa.PrivateKey, error) {
	wif, err := btcutil.DecodeWIF(key.PrivateKey)
	if err != nil {
		return nil, err
	}
	privKey := ecdsa.PrivateKey(*wif.PrivKey)
	return &privKey, nil
}

func (key *bitcoinKey) PriorityCode() uint32 {
	return key.Code
}

func (key *bitcoinKey) Chain() string {
	return key.Network
}

func RandomBitcoinKeyString(chain string) (string, error) {
	priv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", err
	}
	chainParams, err := getChainParams(chain)
	if err != nil {
		return "", err
	}
	wif, err := btcutil.NewWIF(priv, chainParams, false)
	if err != nil {
		return "", err
	}
	return wif.String(), nil
}

func getChainParams(chain string) (*chaincfg.Params, error) {
	switch chain {
	case "regtest":
		return &chaincfg.RegressionNetParams, nil
	case "testnet":
		return &chaincfg.TestNet3Params, nil
	case "mainnet":
		return &chaincfg.MainNetParams, nil
	}
	return nil, fmt.Errorf(ErrPrefix, "Unknown chain")
}
