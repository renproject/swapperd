package fund

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

func (manager *manager) Balances() (map[foundation.TokenName]foundation.Balance, error) {
	balanceMap := map[foundation.TokenName]foundation.Balance{}
	for _, token := range manager.SupportedTokens() {
		balance, err := manager.balance(token.Name)
		if err != nil {
			return balanceMap, err
		}
		balanceMap[token.Name] = balance
	}
	return balanceMap, nil
}

func (manager *manager) balance(token foundation.TokenName) (foundation.Balance, error) {
	switch token {
	case foundation.BTC:
		return manager.balanceBTC()
	case foundation.ETH:
		return manager.balanceETH()
	case foundation.WBTC:
		return manager.balanceERC20(token)
	}
	return foundation.Balance{}, foundation.NewErrUnsupportedToken(token)
}

func (manager *manager) balanceBTC() (foundation.Balance, error) {
	randomKey, err := crypto.GenerateKey()
	if err != nil {
		return foundation.Balance{}, err
	}
	btcAccount := libbtc.NewAccount(libbtc.NewBlockchainInfoClient(manager.config.Bitcoin.Network.Name), randomKey)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	address := manager.config.Bitcoin.Address
	balance, err := btcAccount.Balance(ctx, address, 0)
	if err != nil {
		return foundation.Balance{}, err
	}

	return foundation.Balance{
		Address: address,
		Amount:  big.NewInt(balance).String(),
	}, nil
}

func (manager *manager) balanceETH() (foundation.Balance, error) {
	client, err := beth.Connect(manager.config.Ethereum.Network.URL)
	if err != nil {
		return foundation.Balance{}, err
	}
	address := manager.config.Ethereum.Address

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	balance, err := client.BalanceOf(ctx, common.HexToAddress(address))
	if err != nil {
		return foundation.Balance{}, err
	}

	return foundation.Balance{
		Address: address,
		Amount:  balance.String(),
	}, nil
}

func (manager *manager) balanceERC20(token foundation.TokenName) (foundation.Balance, error) {
	client, err := beth.Connect(manager.config.Ethereum.Network.URL)
	if err != nil {
		return foundation.Balance{}, err
	}
	address := manager.config.Ethereum.Address
	tokenAddr, err := client.ReadAddress(string(token))
	if err != nil {
		return foundation.Balance{}, err
	}
	erc20Contract, err := erc20.NewCompatibleERC20(tokenAddr, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return foundation.Balance{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var balance *big.Int
	if err := client.Get(
		ctx,
		func() error {
			balance, err = erc20Contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
			return err
		},
	); err != nil {
		return foundation.Balance{}, err
	}

	return foundation.Balance{
		Address: address,
		Amount:  balance.String(),
	}, nil
}
