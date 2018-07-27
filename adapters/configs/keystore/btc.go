package keystore

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/ethereum/go-ethereum/crypto"
)

type BitcoinKey struct {
	PrivateKey   *ecdsa.PrivateKey `json:"private_key"`
	CurrencyCode uint32            `json:"currency_code"`
	Network      string            `json:"network"`
}

func NewBitcoinKey(pk string, network string) (BitcoinKey, error) {
	key, err := crypto.HexToECDSA(pk)
	if err != nil {
		return BitcoinKey{}, err
	}
	return BitcoinKey{
		key,
		0,
		network,
	}, nil
}

func (key *BitcoinKey) GetAddress() ([]byte, error) {
	var chainParams *chaincfg.Params

	switch key.Network {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	privKey := (*btcec.PrivateKey)(key.PrivateKey)
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

func (key *BitcoinKey) GetKeyString() (string, error) {
	var chainParams *chaincfg.Params

	switch key.Network {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	privKey := (*btcec.PrivateKey)(key.PrivateKey)
	wif, err := btcutil.NewWIF(privKey, chainParams, false)
	if err != nil {
		return "", err
	}

	return wif.String(), nil
}

func (key *BitcoinKey) GetKey() *ecdsa.PrivateKey {
	return key.PrivateKey
}

func (key *BitcoinKey) PriorityCode() uint32 {
	return key.CurrencyCode
}