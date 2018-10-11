package erc20

import (
	"context"
	"fmt"
	"math/big"

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
	SwapperAddresses map[string]string
}

// NewConnWithConfig creates a new ethereum connection with the given config
// file.
func NewConnWithConfig(config config.EthereumNetwork) (Conn, error) {
	return NewConn(config.URL, config.Network, config.Swappers)
}

// NewConn creates a new ethereum connection with the given config parameters.
func NewConn(url, network, swapperAddresses map[string]string) (Conn, error) {
	ethclient, err := ethclient.Dial(url)
	if err != nil {
		return Conn{}, err
	}
	return Conn{
		Client:           ethclient,
		Network:          network,
		SwapperAddresses: swapperAddresses,
	}, nil
}

// Balance of the given address
func (b *Conn) Balance(address common.Address) (*big.Int, error) {
	binding, err := NewCompatibleERC20(address, bind.ContractBackend(b.Client))
	if err != nil {
		return nil, err
	}
	return binding.BalanceOf(&bind.CallOpts{}, address)
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
		fmt.Println(b.FormatTransactionView("Withdraw transaction successful", tx.Hash().String()))
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

func (b *Conn) FormatTransactionView(msg, txHash string) string {
	switch b.Network {
	case "kovan":
		return fmt.Sprintf("%s, the transaction can be viewed at https://kovan.etherscan.io/tx/%s", msg, txHash)
	case "ropsten":
		return fmt.Sprintf("%s, the transaction can be viewed at https://ropsten.etherscan.io/tx/%s", msg, txHash)
	case "mainnet":
		return fmt.Sprintf("%s, the transaction can be viewed at https://etherscan.io/tx/%s", msg, txHash)
	default:
		panic(fmt.Sprintf("Unknown network :%s", b.Network))
	}
}
