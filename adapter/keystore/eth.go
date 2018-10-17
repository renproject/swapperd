package keystore

import (
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/swapperd/foundation"
)

type EthereumKey struct {
	Network      string
	Address      common.Address
	TransactOpts *bind.TransactOpts
	PrivateKey   *ecdsa.PrivateKey

	*sync.RWMutex
}

func (ethKey EthereumKey) Token() foundation.Token {
	return foundation.TokenWBTC
}

func NewEthereumKey(privKey *ecdsa.PrivateKey, network string) (EthereumKey, error) {
	transactOpts := bind.NewKeyedTransactor(privKey)
	return EthereumKey{
		Network:      network,
		Address:      transactOpts.From,
		TransactOpts: transactOpts,
		PrivateKey:   privKey,
		RWMutex:      new(sync.RWMutex),
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
	return crypto.Sign(crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))), msg), ethKey.PrivateKey)
}

func (ethKey *EthereumKey) SubmitTx(submitTx func(*bind.TransactOpts) error, postCon func() bool) error {
	ethKey.Lock()
	defer ethKey.Unlock()
	for {
		if err := submitTx(ethKey.TransactOpts); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		for i := 0; i < 20; i++ {
			if result := postCon(); result {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}
