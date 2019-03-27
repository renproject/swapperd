package wallet

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

func (wallet *wallet) LockedBalances() map[tokens.Name]*big.Int {
	balancesCopy := map[tokens.Name]*big.Int{}
	for token, balance := range balancesCopy {
		balancesCopy[token] = new(big.Int).SetBytes(balance.Bytes())
	}
	return balancesCopy
}

func (wallet *wallet) LockBalance(token tokens.Name, value string) error {
	val, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return fmt.Errorf("invalid amount: %s", value)
	}
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	wallet.lockedBalances[token] = new(big.Int).Add(wallet.lockedBalances[token], val)
	return nil
}

func (wallet *wallet) UnlockBalance(token tokens.Name, value string) error {
	val, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return fmt.Errorf("invalid amount: %s", value)
	}
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	wallet.lockedBalances[token] = new(big.Int).Sub(wallet.lockedBalances[token], val)
	return nil
}

func (wallet *wallet) AvailableBalance(password string, token tokens.Token) (*big.Int, error) {
	wallet.mu.RLock()
	defer wallet.mu.RUnlock()

	balance, err := wallet.Balance(password, token)
	if err != nil {
		return nil, err
	}
	balanceAmount, ok := new(big.Int).SetString(balance.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to decode balance: %s", balance.Amount)
	}
	return new(big.Int).Add(wallet.lockedBalances[token.Name], balanceAmount), nil
}

func (wallet *wallet) Balance(password string, token tokens.Token) (blockchain.Balance, error) {
	address, err := wallet.GetAddress(password, token.Blockchain)
	if err != nil {
		return blockchain.Balance{}, err
	}

	var balance blockchain.Balance
	switch token.Blockchain {
	case tokens.BITCOIN:
		balance, err = wallet.balanceBTC(address)
	case tokens.ZCASH:
		balance, err = wallet.balanceZEC(address)
	case tokens.ETHEREUM:
		balance, err = wallet.balanceETH(address)
	case tokens.ERC20:
		balance, err = wallet.balanceERC20(token, address)
	default:
		return blockchain.Balance{}, tokens.NewErrUnsupportedBlockchain(token.Blockchain)
	}
	if err != nil {
		return blockchain.Balance{}, err
	}

	val, ok := new(big.Int).SetString(balance.FullAmount, 10)
	if !ok {
		return blockchain.Balance{}, fmt.Errorf("invalid balance: %d", val)
	}
	balance.Amount = new(big.Int).Sub(val, wallet.lockedBalances[token.Name]).String()
	return balance, nil
}

func (wallet *wallet) balanceBTC(address string) (blockchain.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	btcClient, err := wallet.bitcoinClient()
	if err != nil {
		return blockchain.Balance{}, err
	}

	balance, err := btcClient.Balance(ctx, address, 0)
	if err != nil {
		return blockchain.Balance{}, err
	}

	return blockchain.Balance{
		Address:    address,
		Decimals:   int(tokens.BTC.Decimals),
		FullAmount: big.NewInt(balance).String(),
	}, nil
}

func (wallet *wallet) balanceZEC(address string) (blockchain.Balance, error) {
	zecClient, err := wallet.zcashClient()
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
		Address:    address,
		Decimals:   int(tokens.ZEC.Decimals),
		FullAmount: big.NewInt(balance).String(),
	}, nil
}

func (wallet *wallet) balanceETH(address string) (blockchain.Balance, error) {
	client, err := wallet.ethereumClient()
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
		Address:    address,
		Decimals:   int(tokens.ETH.Decimals),
		FullAmount: balance.String(),
	}, nil
}

func (wallet *wallet) balanceERC20(token tokens.Token, address string) (blockchain.Balance, error) {
	client, err := wallet.ethereumClient()
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
		Address:    address,
		Decimals:   int(token.Decimals),
		FullAmount: balance.String(),
	}, nil
}
