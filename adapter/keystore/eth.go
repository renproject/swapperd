package keystore

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

type EthereumKey struct {
	Network      string
	Address      common.Address
	TransactOpts *bind.TransactOpts
	PrivateKey   *ecdsa.PrivateKey
}

func (ethKey EthereumKey) Token() token.Token {
	return token.ETH
}

func NewEthereumKey(privKey *ecdsa.PrivateKey, network string) (EthereumKey, error) {
	transactOpts := bind.NewKeyedTransactor(privKey)
	return EthereumKey{
		Network:      network,
		Address:      transactOpts.From,
		TransactOpts: transactOpts,
		PrivateKey:   privKey,
	}, nil
}

func RandomEthereumKey(network string) (EthereumKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return EthereumKey{}, nil
	}
	return NewEthereumKey(privKey, network)
}

func (ethKey *EthereumKey) Sign(msg []byte) ([]byte, error) {
	return crypto.Sign(msg, ethKey.PrivateKey)
}
