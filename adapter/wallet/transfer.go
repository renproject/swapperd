package wallet

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error) {
	switch token {
	case blockchain.TokenBTC:
		return wallet.transferBTC(password, to, amount)
	case blockchain.TokenETH:
		return wallet.transferETH(password, to, amount)
	case blockchain.TokenWBTC:
		return wallet.transferERC20(password, token, to, amount)
	}
	return "", blockchain.NewErrUnsupportedToken(token.Name)
}

func (wallet *wallet) transferBTC(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	account, err := wallet.BitcoinAccount(password)
	if err != nil {
		return "", err
	}
	return "", account.Transfer(ctx, to, amount.Int64())
}

func (wallet *wallet) transferETH(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return "", err
	}
	return "", account.Transfer(ctx, common.HexToAddress(to), amount, 1)
}

func (wallet *wallet) transferERC20(password string, token blockchain.Token, to string, amount *big.Int) (string, error) {
	var txHash string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return txHash, err
	}
	tokenAddress, err := account.ReadAddress(string(token.Name))
	if err != nil {
		return txHash, err
	}

	tokenContract, err := erc20.NewCompatibleERC20(tokenAddress, bind.ContractBackend(account.EthClient()))
	if err != nil {
		return txHash, err
	}

	if err := account.Transact(
		ctx,
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := tokenContract.Transfer(tops, common.HexToAddress(to), amount)
			if err != nil {
				return tx, err
			}
			txHash = tx.Hash().String()
			return tx, nil
		},
		nil,
		1,
	); err != nil {
		return txHash, err
	}

	return txHash, nil
}
