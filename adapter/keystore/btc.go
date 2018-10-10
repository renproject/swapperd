package keystore

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/republicprotocol/swapperd/utils"
)

type BitcoinKey struct {
	Network       *chaincfg.Params
	AddressString string
	Address       btcutil.Address
	PKHash        *btcutil.AddressPubKeyHash
	WIF           *btcutil.WIF
	PrivateKey    *btcec.PrivateKey
	PublicKey     []byte
	Compressed    bool
}

func (btcKey BitcoinKey) Token() foundation.Token {
	return foundation.TokenBTC
}

func NewBitcoinKey(wifString string, network string) (BitcoinKey, error) {
	net := utils.GetChainParams(network)
	wif, err := btcutil.DecodeWIF(wifString)
	if err != nil {
		return BitcoinKey{}, err
	}

	var pubKeyBytes []byte
	var compressed bool
	if network == "mainnet" {
		pubKeyBytes = wif.PrivKey.PubKey().SerializeCompressed()
		compressed = true
	} else {
		pubKeyBytes = wif.PrivKey.PubKey().SerializeUncompressed()
		compressed = false
	}

	pubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, net)
	if err != nil {
		return BitcoinKey{}, err
	}
	addrString := pubKey.EncodeAddress()
	addr, err := btcutil.DecodeAddress(addrString, net)
	if err != nil {
		return BitcoinKey{}, err
	}

	return BitcoinKey{
		Network:       net,
		WIF:           wif,
		AddressString: addrString,
		Address:       addr,
		PrivateKey:    wif.PrivKey,
		PKHash:        pubKey.AddressPubKeyHash(),
		PublicKey:     pubKeyBytes,
		Compressed:    compressed,
	}, nil

}

func RandomBitcoinKey(network string) (BitcoinKey, error) {
	key, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return BitcoinKey{}, err
	}
	wif, err := btcutil.NewWIF(key, utils.GetChainParams(network), true)
	if err != nil {
		return BitcoinKey{}, err
	}
	return NewBitcoinKey(wif.String(), network)
}
