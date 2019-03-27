package wallet

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/renproject/libbtc-go"
	"github.com/renproject/libeth-go"
	"github.com/renproject/libzec-go"
	"github.com/renproject/swapperd/adapter/binder/erc20"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/tokens"
	"github.com/republicprotocol/co-go"
)

func (wallet *wallet) Balances(password string) (map[tokens.Name]blockchain.Balance, error) {
	balanceMap := map[tokens.Name]blockchain.Balance{}
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

func (wallet *wallet) Balance(password string, token tokens.Token) (blockchain.Balance, error) {
	address, err := wallet.GetAddress(password, token.Blockchain)
	if err != nil {
		return blockchain.Balance{}, err
	}

	switch token.Blockchain {
	case tokens.BITCOIN:
		return wallet.balanceBTC(address)
	case tokens.ZCASH:
		return wallet.balanceZEC(address)
	case tokens.ETHEREUM:
		return wallet.balanceETH(address)
	case tokens.ERC20:
		return wallet.balanceERC20(token, address)
	default:
		return blockchain.Balance{}, tokens.NewErrUnsupportedBlockchain(token.Blockchain)
	}
}

func (wallet *wallet) balanceBTC(address string) (blockchain.Balance, error) {
	btcClient, err := libbtc.NewBlockchainInfoClient(wallet.config.Bitcoin.Network.Name)
	if err != nil {
		return blockchain.Balance{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	balance, err := btcClient.Balance(ctx, address, 0)
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address:  address,
		Decimals: int(tokens.BTC.Decimals),
		Amount:   big.NewInt(balance).String(),
	}, nil
}

func (wallet *wallet) balanceZEC(address string) (blockchain.Balance, error) {
	zecClient, err := libzec.NewMercuryClient(wallet.config.ZCash.Network.Name)
	if err != nil {
		return blockchain.Balance{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	balance, err := zecClient.Balance(ctx, address, 0)
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address:  address,
		Decimals: int(tokens.ZEC.Decimals),
		Amount:   big.NewInt(balance).String(),
	}, nil
}

func (wallet *wallet) balanceETH(address string) (blockchain.Balance, error) {
	client, err := libeth.NewInfuraClient(wallet.config.Ethereum.Network.Name, "172978c53e244bd78388e6d50a4ae2fa")
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
		Decimals: int(tokens.ETH.Decimals),
		Amount:   balance.String(),
	}, nil
}

func (wallet *wallet) balanceERC20(token tokens.Token, address string) (blockchain.Balance, error) {
	client, err := libeth.NewInfuraClient(wallet.config.Ethereum.Network.Name, "172978c53e244bd78388e6d50a4ae2fa")
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
		Decimals: int(token.Decimals),
		Amount:   balance.String(),
	}, nil
}
