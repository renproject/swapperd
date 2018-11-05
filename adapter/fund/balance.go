package fund

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

func (manager *manager) Balances() (map[foundation.Token]Balance, error) {
	balanceMap := map[foundation.Token]Balance{}
	for _, token := range manager.supportedTokens {
		balance, err := manager.balance(token)
		if err != nil {
			return balanceMap, err
		}
		balanceMap[token] = balance
	}
	return balanceMap, nil
}

func (manager *manager) balance(token foundation.Token) (Balance, error) {
	switch token {
	case foundation.TokenBTC:
		return manager.balanceBTC()
	case foundation.TokenETH:
		return manager.balanceETH()
	case foundation.TokenWBTC:
		return manager.balanceERC20(token)
	}
	return Balance{}, foundation.NewErrUnsupportedToken(token.Name)
}

func (manager *manager) balanceBTC() (Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	address := manager.addresses[foundation.Bitcoin]
	balance, err := manager.btcAccount.Balance(ctx, address, 0)
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: address,
		Amount:  big.NewInt(balance),
	}, nil
}

func (manager *manager) balanceETH() (Balance, error) {
	client := manager.ethAccount.EthClient()
	address := manager.addresses[foundation.Ethereum]
	balance, err := client.BalanceOf(context.Background(), common.HexToAddress(address))
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: address,
		Amount:  balance,
	}, nil
}

func (manager *manager) balanceERC20(token foundation.Token) (Balance, error) {
	client := manager.ethAccount.EthClient()
	address := manager.addresses[foundation.Ethereum]

	tokenAddr, err := manager.ethAccount.ReadAddress(fmt.Sprintf("ERC20:%s", token.Name))
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
			balance, err = erc20Contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
			return err
		},
	); err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: address,
		Amount:  balance,
	}, nil
}
