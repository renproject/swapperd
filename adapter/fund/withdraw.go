package fund

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

func (manager *manager) Withdraw(password string, token foundation.Token, to string, amount *big.Int) (string, error) {
	switch token {
	case foundation.TokenBTC:
		return manager.withdrawBTC(password, to, amount)
	case foundation.TokenETH:
		return manager.withdrawETH(password, to, amount)
	case foundation.TokenWBTC:
		return manager.withdrawERC20(password, token, to, amount)
	}
	return "", foundation.NewErrUnsupportedToken(token.Name)
}

func (manager *manager) withdrawBTC(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if amount == nil {
		balance, err := manager.balanceBTC()
		if err != nil {
			return "", err
		}
		amount = balance.Amount
	}
	account, err := manager.BitcoinAccount(password)
	if err != nil {
		return "", err
	}
	return "", account.Transfer(ctx, to, amount.Int64())
}

func (manager *manager) withdrawETH(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if amount == nil {
		balance, err := manager.balanceETH()
		if err != nil {
			return "", err
		}
		amount = balance.Amount
	}
	account, err := manager.EthereumAccount(password)
	if err != nil {
		return "", err
	}
	return "", account.Transfer(ctx, common.HexToAddress(to), amount, 1)
}

func (manager *manager) withdrawERC20(password string, token foundation.Token, to string, amount *big.Int) (string, error) {
	var txHash string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if amount == nil {
		balance, err := manager.balanceERC20(token)
		if err != nil {
			return txHash, err
		}
		amount = balance.Amount
	}
	account, err := manager.EthereumAccount(password)
	if err != nil {
		return txHash, err
	}
	tokenAddress, err := account.ReadAddress(fmt.Sprintf("ERC20:%s", token.Name))
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
