package balance

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/republicprotocol/swapperd/adapter/account"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

type Query struct {
	password string
	response chan<- map[foundation.Token]Balance
	err      chan<- error
}

type Balance struct {
	Address string
	Amount  *big.Int
}

func NewQuery(password string) (Query, <-chan map[foundation.Token]Balance, <-chan error) {
	response := make(chan map[foundation.Token]Balance, 1)
	err := make(chan error, 1)
	return Query{
		password: password,
		response: response,
		err:      err,
	}, response, err
}

type Book interface {
	Run(done <-chan struct{}, queries <-chan Query)
}

type book struct {
	accounts        account.Accounts
	supportedTokens []foundation.Token
}

func New(accounts account.Accounts) Book {
	return &book{
		accounts: accounts,
		supportedTokens: []foundation.Token{
			foundation.TokenBTC,
			foundation.TokenETH,
			foundation.TokenWBTC,
		},
	}
}

func (book *book) Run(done <-chan struct{}, queries <-chan Query) {
	for {
		select {
		case <-done:
			return
		case query, ok := <-queries:
			if !ok {
				return
			}
			balanceMap, err := book.balances(query.password)
			query.err <- err
			query.response <- balanceMap
		}
	}
}

func (book *book) balances(password string) (map[foundation.Token]Balance, error) {
	balanceMap := map[foundation.Token]Balance{}
	for _, token := range book.supportedTokens {
		balance, err := book.getBalance(token, password)
		if err != nil {
			return balanceMap, err
		}
		balanceMap[token] = balance
	}
	return balanceMap, nil
}

func (book *book) getBalance(token foundation.Token, password string) (Balance, error) {
	switch token {
	case foundation.TokenBTC:
		return book.btcBalance(password)
	case foundation.TokenETH:
		return book.ethBalance(password)
	case foundation.TokenWBTC:
		return book.erc20Balance(password, token)
	}
	return Balance{}, foundation.NewErrUnsupportedToken(token.Name)
}

func (book *book) btcBalance(password string) (Balance, error) {
	account, err := book.accounts.GetBitcoinAccount(password)
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

func (book *book) ethBalance(password string) (Balance, error) {
	account, err := book.accounts.GetEthereumAccount(password)
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

func (book *book) erc20Balance(password string, token foundation.Token) (Balance, error) {
	account, err := book.accounts.GetEthereumAccount(password)
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
