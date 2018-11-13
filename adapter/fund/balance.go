package fund

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/beth-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

func (manager *manager) Balances() (map[foundation.Token]Balance, error) {
	balanceMap := map[foundation.Token]Balance{}
	for _, token := range manager.SupportedTokens() {
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
	randomKey, err := crypto.GenerateKey()
	if err != nil {
		return Balance{}, err
	}
	btcAccount := libbtc.NewAccount(libbtc.NewBlockchainInfoClient(manager.config.Bitcoin.Network.Name), randomKey)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	address := manager.config.Bitcoin.Address
	balance, err := btcAccount.Balance(ctx, address, 0)
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		Address: address,
		Amount:  big.NewInt(balance),
	}, nil
}

func (manager *manager) balanceETH() (Balance, error) {
	client, err := beth.Connect(manager.config.Ethereum.Network.URL)
	if err != nil {
		return Balance{}, err
	}
	address := manager.config.Ethereum.Address
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
	client, err := beth.Connect(manager.config.Ethereum.Network.URL)
	if err != nil {
		return Balance{}, err
	}
	address := manager.config.Ethereum.Address
	tokenAddr, err := client.ReadAddress(token.Name)
	if err != nil {
		return Balance{}, err
	}
	erc20Contract, err := erc20.NewCompatibleERC20(tokenAddr, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return Balance{}, err
	}
	fmt.Println(tokenAddr.String())
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
