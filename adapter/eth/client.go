package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

type Conn struct {
	Network            string
	Client             *ethclient.Client
	RenExAtomicSwapper common.Address
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
		Client:             ethclient,
		Network:            network,
		RenExAtomicSwapper: common.HexToAddress(swapperAddress),
	}, nil
}

// Balance of the given address
func (b *Conn) Balance(address common.Address) (*big.Int, error) {
	return b.Client.PendingBalanceAt(context.Background(), address)
}

// Transfer is a helper function for sending ETH to an address
func (b *Conn) Transfer(to common.Address, key keystore.EthereumKey, value *big.Int) error {
	if value.Cmp(big.NewInt(0)) == 0 {
		balance, err := b.Balance(to)
		if err != nil {
			return err
		}
		value = balance.Sub(balance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(14), nil))
	}

	nonceBefore, err := b.Client.PendingNonceAt(context.Background(), key.Address)
	if err != nil {
		return err
	}

	key.SubmitTx(func(tops *bind.TransactOpts) error {
		tops.Value = value
		// Why is there no ethclient.Transfer?
		bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.Client, nil)
		_, err := bound.Transfer(tops)
		if err != nil {
			return err
		}
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
