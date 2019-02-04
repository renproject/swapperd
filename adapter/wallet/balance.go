package wallet

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/renproject/swapperd/adapter/binder/erc20"
	"github.com/renproject/swapperd/foundation/blockchain"
)

func (wallet *wallet) Balances(password string) (map[blockchain.TokenName]blockchain.Balance, error) {
	balanceMap := map[blockchain.TokenName]blockchain.Balance{}
	mu := new(sync.RWMutex)
	tokens := wallet.SupportedTokens()
	co.ParForAll(tokens, func(i int) {
		token := tokens[i]
		balance, err := wallet.Balance(password, token)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		balanceMap[token.Name] = balance
	})
	return balanceMap, nil
}

func (wallet *wallet) Balance(password string, token blockchain.Token) (blockchain.Balance, error) {
	address, err := wallet.GetAddress(password, token.Blockchain)
	if err != nil {
		return blockchain.Balance{}, err
	}

	switch token.Blockchain {
	case blockchain.Bitcoin:
		return wallet.balanceBTC(address)
	case blockchain.Ethereum:
		return wallet.balanceETH(address)
	case blockchain.ERC20:
		return wallet.balanceERC20(token, address)
	default:
		return blockchain.Balance{}, blockchain.NewErrUnsupportedToken(token.Name)
	}
}

func (wallet *wallet) balanceBTC(address string) (blockchain.Balance, error) {
	randomKey, err := crypto.GenerateKey()
	if err != nil {
		return blockchain.Balance{}, err
	}
	btcAccount := libbtc.NewAccount(libbtc.NewBlockchainInfoClient(wallet.config.Bitcoin.Network.Name), randomKey)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	balance, err := btcAccount.Balance(ctx, address, 0)
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address:  address,
		Decimals: blockchain.TokenBTC.Decimals,
		Amount:   big.NewInt(balance).String(),
	}, nil
}

func (wallet *wallet) balanceETH(address string) (blockchain.Balance, error) {
	client, err := beth.Connect(wallet.config.Ethereum.Network.URL)
	if err != nil {
		return blockchain.Balance{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	balance, err := client.BalanceOf(ctx, common.HexToAddress(address))
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address:  address,
		Decimals: blockchain.TokenETH.Decimals,
		Amount:   balance.String(),
	}, nil
}

func (wallet *wallet) balanceERC20(token blockchain.Token, address string) (blockchain.Balance, error) {
	client, err := beth.Connect(wallet.config.Ethereum.Network.URL)
	if err != nil {
		return blockchain.Balance{}, err
	}
	tokenAddr, err := client.ReadAddress(string(token.Name))
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
		Address:  address,
		Decimals: token.Decimals,
		Amount:   balance.String(),
	}, nil
}
