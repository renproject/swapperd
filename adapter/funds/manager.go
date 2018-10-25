package funds

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/swapperd/adapter/account"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

type Balance struct {
	Address string
	Amount  *big.Int
}

type Manager interface {
	SupportedTokens() []foundation.Token
	Withdraw(password string, token foundation.Token, to string, amount *big.Int) (string, error)
	Balances(password string) (map[foundation.Token]Balance, error)
}

type manager struct {
	accounts        account.Accounts
	supportedTokens []foundation.Token
}

func New(accounts account.Accounts) Manager {
	return &manager{
		accounts: accounts,
		supportedTokens: []foundation.Token{
			foundation.TokenBTC,
			foundation.TokenETH,
			foundation.TokenWBTC,
		},
	}
}

func (manager *manager) SupportedTokens() []foundation.Token {
	return manager.supportedTokens
}

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

func (manager *manager) Balances(password string) (map[foundation.Token]Balance, error) {
	balanceMap := map[foundation.Token]Balance{}
	for _, token := range manager.supportedTokens {
		balance, err := manager.balance(token, password)
		if err != nil {
			return balanceMap, err
		}
		balanceMap[token] = balance
	}
	return balanceMap, nil
}

func (manager *manager) balance(token foundation.Token, password string) (Balance, error) {
	switch token {
	case foundation.TokenBTC:
		return manager.balanceBTC(password)
	case foundation.TokenETH:
		return manager.balanceETH(password)
	case foundation.TokenWBTC:
		return manager.balanceERC20(password, token)
	}
	return Balance{}, foundation.NewErrUnsupportedToken(token.Name)
}

func (manager *manager) balanceBTC(password string) (Balance, error) {
	account, err := manager.accounts.GetBitcoinAccount(password)
	if err != nil {
		return Balance{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	address, err := account.Address()
	if err != nil {
		return Balance{}, err
	}
	balance, err := account.Balance(ctx, address.EncodeAddress(), 0)
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: address.EncodeAddress(),
		Amount:  big.NewInt(balance),
	}, nil
}

func (manager *manager) withdrawBTC(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if amount == nil {
		balance, err := manager.balanceBTC(password)
		if err != nil {
			return "", err
		}
		amount = balance.Amount
	}
	account, err := manager.accounts.GetBitcoinAccount(password)
	if err != nil {
		return "", err
	}
	return "", account.Transfer(ctx, to, amount.Int64())
}

func (manager *manager) balanceETH(password string) (Balance, error) {
	account, err := manager.accounts.GetEthereumAccount(password)
	if err != nil {
		return Balance{}, err
	}
	client := account.EthClient()
	balance, err := client.BalanceOf(context.Background(), account.Address())
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: account.Address().String(),
		Amount:  balance,
	}, nil
}

func (manager *manager) withdrawETH(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if amount == nil {
		balance, err := manager.balanceETH(password)
		if err != nil {
			return "", err
		}
		amount = balance.Amount
	}
	account, err := manager.accounts.GetEthereumAccount(password)
	if err != nil {
		return "", err
	}
	return "", account.Transfer(ctx, common.HexToAddress(to), amount, 1)
}

func (manager *manager) balanceERC20(password string, token foundation.Token) (Balance, error) {
	account, err := manager.accounts.GetEthereumAccount(password)
	if err != nil {
		return Balance{}, err
	}
	client := account.EthClient()

	tokenAddr, err := account.ReadAddress(fmt.Sprintf("ERC20:%s", token.Name))
	if err != nil {
		return Balance{}, err
	}
	erc20Contract, err := erc20.NewCompatibleERC20(tokenAddr, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return Balance{}, err
	}
	var balance *big.Int
	if err := client.Get(
		context.Background(),
		func() error {
			balance, err = erc20Contract.BalanceOf(&bind.CallOpts{}, account.Address())
			return err
		},
	); err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: account.Address().String(),
		Amount:  balance,
	}, nil
}

func (manager *manager) withdrawERC20(password string, token foundation.Token, to string, amount *big.Int) (string, error) {
	var txHash string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if amount == nil {
		balance, err := manager.balanceERC20(password, token)
		if err != nil {
			return txHash, err
		}
		amount = balance.Amount
	}
	account, err := manager.accounts.GetEthereumAccount(password)
	if err != nil {
		return txHash, err
	}
	tokenAddress, err := account.ReadAddress(fmt.Sprintf("ERC20:%s", token.Name))
	if err != nil {
		return txHash, err
	}

	ethClient := account.EthClient()
	tokenContract, err := erc20.NewCompatibleERC20(tokenAddress, bind.ContractBackend(ethClient.EthClient()))
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
