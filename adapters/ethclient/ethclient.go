package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Network is used to represent an Ethereum chain
type Network string

const (
	// NetworkMainnet represents the Ethereum mainnet
	NetworkMainnet Network = "mainnet"
	// NetworkRopsten represents the Ethereum Ropsten testnet
	NetworkRopsten Network = "ropsten"
	// NetworkKovan represents the Ethereum Kovan testnet
	NetworkKovan Network = "kovan"
	// NetworkGanache represents a Ganache testrpc server
	NetworkGanache Network = "ganache"
)

// Connection contains the client and the contracts deployed to it
type Connection struct {
	Client         *ethclient.Client
	EthAddress     common.Address
	NetworkAddress common.Address
	Network        Network
}

// Connect to a URI.
func Connect(network Network) (Connection, error) {

	var uri string
	var ethSwapAddress string

	switch network {
	case NetworkGanache:
		uri = "http://localhost:8545"
	case NetworkRopsten:
		uri = "https://ropsten.infura.io"
		ethSwapAddress = ""
	case NetworkKovan:
		uri = "https://kovan.infura.io"
		ethSwapAddress = ""
	default:
		return Connection{}, fmt.Errorf("cannot connect to %s: unsupported", network)
	}

	ethclient, err := ethclient.Dial(uri)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		Client:     ethclient,
		EthAddress: common.HexToAddress(ethSwapAddress),
		Network:    network,
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Connection) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	switch b.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	default:
		return bind.WaitMined(ctx, b.Client, tx)
	}
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Connection) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	switch b.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return common.Address{}, nil
	default:
		return bind.WaitDeployed(ctx, b.Client, tx)
	}
}

// TransferEth is a helper function for sending ETH to an address
func (b *Connection) TransferEth(ctx context.Context, from *bind.TransactOpts, to common.Address, value *big.Int) error {
	transactor := &bind.TransactOpts{
		From:     from.From,
		Nonce:    from.Nonce,
		Signer:   from.Signer,
		Value:    value,
		GasPrice: from.GasPrice,
		GasLimit: 30000,
		Context:  from.Context,
	}

	// Why is there no ethclient.Transfer?
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.Client, nil)
	tx, err := bound.Transfer(transactor)
	if err != nil {
		return err
	}
	_, err = b.PatchedWaitMined(ctx, tx)
	return err
}
