package ethclient

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/republicprotocol/atom-go/adapters/configs/network"
)

type Conn struct {
	network          string
	client           *ethclient.Client
	atomAddress      common.Address
	walletAddress    common.Address
	infoAddress      common.Address
	orderBookAddress common.Address
}

// Connect to an ethereum network.
func Connect(config network.Config) (Conn, error) {
	ethclient, err := ethclient.Dial(config.Ethereum.URL)
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		client:           ethclient,
		network:          config.Ethereum.Chain,
		atomAddress:      common.HexToAddress(config.Ethereum.AtomAddress),
		infoAddress:      common.HexToAddress(config.Ethereum.InfoAddress),
		walletAddress:    common.HexToAddress(config.Ethereum.WalletAddress),
		orderBookAddress: common.HexToAddress(config.Ethereum.OrderBookAddress),
	}, nil
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
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.client, nil)
	tx, err := bound.Transfer(transactor)
	if err != nil {
		return err
	}
	_, err = b.PatchedWaitMined(context.Background(), tx)
	return err
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Conn) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	switch b.network {
	case "ganache":
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	default:
		return bind.WaitMined(ctx, b.client, tx)
	}
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Conn) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	switch b.network {
	case "ganache":
		time.Sleep(100 * time.Millisecond)
		return common.Address{}, nil
	default:
		return bind.WaitDeployed(ctx, b.client, tx)
	}
}

func (conn *Conn) AtomAddress() common.Address {
	return conn.atomAddress
}

func (conn *Conn) WalletAddress() common.Address {
	return conn.walletAddress
}

func (conn *Conn) InfoAddress() common.Address {
	return conn.infoAddress
}

func (conn *Conn) OrderBookAddress() common.Address {
	return conn.orderBookAddress
}

func (conn *Conn) Network() string {
	return conn.network
}

func (conn *Conn) Client() *ethclient.Client {
	return conn.client
}
