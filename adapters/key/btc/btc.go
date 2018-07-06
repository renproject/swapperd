package btc

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/atom-go/services/swap"
)

type bitcoinKey struct {
	privateKey *ecdsa.PrivateKey
	network    string
}

func NewBitcoinKey(pk string, network string) (swap.Key, error) {
	key, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}
	return &bitcoinKey{
		key,
		network,
	}, nil
}

func (key bitcoinKey) GetAddress() ([]byte, error) {
	var chainParams *chaincfg.Params

	switch key.network {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	privKey := (*btcec.PrivateKey)(key.privateKey)
	wif, err := btcutil.NewWIF(privKey, chainParams, false)
	if err != nil {
		return nil, err
	}

	spubKey := wif.SerializePubKey()
	pubKey, err := btcutil.NewAddressPubKey(spubKey, chainParams)
	if err != nil {
		return nil, err
	}

	return []byte(pubKey.EncodeAddress()), nil
}

func (key bitcoinKey) GetKeyString() (string, error) {
	var chainParams *chaincfg.Params

	switch key.network {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	privKey := (*btcec.PrivateKey)(key.privateKey)
	wif, err := btcutil.NewWIF(privKey, chainParams, false)
	if err != nil {
		return "", err
	}

	return wif.String(), nil
}

func (key bitcoinKey) GetKey() *ecdsa.PrivateKey {
	return key.privateKey
}

func (key bitcoinKey) PriorityCode() uint32 {
	return 0
}
