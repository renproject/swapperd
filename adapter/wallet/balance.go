package wallet

import (
	"context"
	"math/big"
	"time"

	"github.com/republicprotocol/beth-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) Balances() (map[blockchain.TokenName]blockchain.Balance, error) {
	balanceMap := map[blockchain.TokenName]blockchain.Balance{}
	for _, token := range wallet.SupportedTokens() {
		balance, err := wallet.balance(token.Name)
		if err != nil {
			return balanceMap, err
		}
		balanceMap[token.Name] = balance
	}
	return balanceMap, nil
}

func (wallet *wallet) balance(token blockchain.TokenName) (blockchain.Balance, error) {
	switch token {
	case blockchain.BTC:
		return wallet.balanceBTC()
	case blockchain.ETH:
		return wallet.balanceETH()
	case blockchain.WBTC:
		return wallet.balanceERC20(token)
	}
	return blockchain.Balance{}, blockchain.NewErrUnsupportedToken(token)
}

func (wallet *wallet) balanceBTC() (blockchain.Balance, error) {
	randomKey, err := crypto.GenerateKey()
	if err != nil {
		return blockchain.Balance{}, err
	}
	btcAccount := libbtc.NewAccount(libbtc.NewBlockchainInfoClient(wallet.config.Bitcoin.Network.Name), randomKey)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	address := wallet.config.Bitcoin.Address
	balance, err := btcAccount.Balance(ctx, address, 0)
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address: address,
		Amount:  big.NewInt(balance).String(),
	}, nil
}

func (wallet *wallet) balanceETH() (blockchain.Balance, error) {
	client, err := beth.Connect(wallet.config.Ethereum.Network.URL)
	if err != nil {
		return blockchain.Balance{}, err
	}
	address := wallet.config.Ethereum.Address

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	balance, err := client.BalanceOf(ctx, common.HexToAddress(address))
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address: address,
		Amount:  balance.String(),
	}, nil
}

func (wallet *wallet) balanceERC20(token blockchain.TokenName) (blockchain.Balance, error) {
	client, err := beth.Connect(wallet.config.Ethereum.Network.URL)
	if err != nil {
		return blockchain.Balance{}, err
	}
	address := wallet.config.Ethereum.Address
	tokenAddr, err := client.ReadAddress(string(token))
	if err != nil {
		return blockchain.Balance{}, err
	}
	erc20Contract, err := erc20.NewCompatibleERC20(tokenAddr, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return blockchain.Balance{}, err
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
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address: address,
		Amount:  balance.String(),
	}, nil
}
