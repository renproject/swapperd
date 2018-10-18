package balance

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/swapperd/adapter/account"
	"github.com/republicprotocol/swapperd/adapter/eth/erc20"
	"github.com/republicprotocol/swapperd/foundation"
)

type Request struct {
	token    foundation.Token
	response chan<- *big.Int
	err      chan<- error
}

func GetBalance(accounts account.Accounts, requestCh <-chan Request) {
	for {
		select {
		case req, ok := <-requestCh:
			if !ok {
				return
			}
			switch req.token {
			case foundation.TokenBTC:
				account := accounts.GetBitcoinAccount()
				address, err := account.Address()
				if err != nil {
					req.err <- err
					continue
				}
				balance, err := account.Balance(context.Background(), address.String(), 0)
				if err != nil {
					req.err <- err
					continue
				}
				req.response <- big.NewInt(balance)
			case foundation.TokenETH:
				account := accounts.GetEthereumAccount()
				client := account.EthClient()
				bal, err := client.BalanceOf(context.Background(), account.Address())
				if err != nil {
					req.err <- err
					continue
				}
				req.response <- bal
			case foundation.TokenWBTC:
				bal, err := erc20Balance(accounts.GetEthereumAccount(), foundation.TokenWBTC)
				if err != nil {
					req.err <- err
					continue
				}
				req.response <- bal
			default:
				req.err <- foundation.NewErrUnsupportedToken(req.token.Name)
			}
		}
	}
}

func erc20Balance(account beth.Account, token foundation.Token) (balance *big.Int, err error) {
	client := account.EthClient()
	tokenAddr, err := account.ReadAddress(fmt.Sprintf("ERC20:%s", token.Name))
	if err != nil {
		return
	}
	erc20Contract, err := erc20.NewCompatibleERC20(tokenAddr, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return
	}
	err = client.Get(
		context.Background(),
		func() error {
			balance, err = erc20Contract.BalanceOf(&bind.CallOpts{}, account.Address())
			return err
		},
	)
	return
}
