package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
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

// NewAccount creates a new account and funds it with ether
func (b *Conn) NewAccount(value int64, from *bind.TransactOpts) (common.Address, *bind.TransactOpts, error) {
	account, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, &bind.TransactOpts{}, err
	}

	accountAddress := crypto.PubkeyToAddress(account.PublicKey)
	accountAuth := bind.NewKeyedTransactor(account)

	return accountAddress, accountAuth, b.Transfer(accountAddress, from, value)
}

// Transfer is a helper function for sending ETH to an address
func (b *Conn) Transfer(to common.Address, from *bind.TransactOpts, value int64) error {
	transactor := &bind.TransactOpts{
		From:     from.From,
		Nonce:    from.Nonce,
		Signer:   from.Signer,
		Value:    big.NewInt(value),
		GasPrice: from.GasPrice,
		GasLimit: 30000,
		Context:  from.Context,
	}

	// Why is there no ethclient.Transfer?
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.Client, nil)
	_, err := bound.Transfer(transactor)
	if err != nil {
		return err
	}
	return nil
}
