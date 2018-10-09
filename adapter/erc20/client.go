package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/domain/token"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/keystore"
)

type Conn struct {
	Network          string
	Client           *ethclient.Client
	SwapperAddresses map[token.Token]common.Address
}

// NewConnWithConfig creates a new ethereum connection with the given config
// file.
func NewConnWithConfig(config config.EthereumNetwork) (Conn, error) {
	return NewConn(config.URL, config.Network, config.Swapper)
}

// NewConn creates a new ethereum connection with the given config parameters.
func NewConn(url, network, swapperAddress string) (Conn, error) {
	ethclient, err := ethclient.Dial(url)
	if err != nil {
		return Conn{}, err
	}
	return Conn{
		Client:           ethclient,
		Network:          network,
		SwapperAddresses: common.HexToAddress(swapperAddress),
	}, nil
}

// Balance of the given address
func (b *Conn) Balance(address common.Address) (*big.Int, error) {
	return b.Client.PendingBalanceAt(context.Background(), address)
}

// Transfer is a helper function for sending ETH to an address
func (b *Conn) Transfer(to common.Address, key keystore.EthereumKey, value *big.Int) error {
	balance, err := b.Balance(key.Address)
	if err != nil {
		return err
	}

	if value.Cmp(balance) > 0 {
		return fmt.Errorf("Not enough balance expected withdrawal: %v balance: %v", value, balance)
	}

	fee := big.NewInt(0).Mul(big.NewInt(21), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(13), nil))
	if balance.Cmp(fee) <= 0 {
		return fmt.Errorf("Not enough balance: %v", balance)
	}

	if value.Cmp(big.NewInt(0)) == 0 {
		value = balance.Sub(balance, fee)
	}

	nonceBefore, err := b.Client.PendingNonceAt(context.Background(), key.Address)
	if err != nil {
		return err
	}
	key.SubmitTx(func(tops *bind.TransactOpts) error {
		txOpts := &bind.TransactOpts{
			From:     tops.From,
			Nonce:    tops.Nonce,
			Signer:   tops.Signer,
			Value:    value,
			GasPrice: big.NewInt(10000000000),
			GasLimit: 21000,
			Context:  tops.Context,
		}
		// Why is there no ethclient.Transfer?
		bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.Client, nil)
		tx, err := bound.Transfer(txOpts)
		if err != nil {
			return err
		}
		fmt.Printf("Transaction can be viewed at https://etherscan.io/tx/%s\n", tx.Hash().String())
		return nil
	}, func() bool {
		nonceAfter, err := b.Client.PendingNonceAt(context.Background(), key.Address)
		if err != nil {
			return false
		}
		return nonceAfter > nonceBefore
	},
	)
	return nil
}
