package account

import (
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
)

type Accounts interface {
	GetBitcoinAccount() libbtc.Account
	GetEthereumAccount() beth.Account
}

type accounts struct {
	bitcoinAccount  libbtc.Account
	ethereumAccount beth.Account
}

func New(btc libbtc.Account, eth beth.Account) Accounts {
	return &accounts{
		bitcoinAccount:  btc,
		ethereumAccount: eth,
	}
}

// GetBitcoinAccount returns the bitcoin account
func (accounts *accounts) GetBitcoinAccount() libbtc.Account {
	return accounts.bitcoinAccount
}

// GetKey returns the key object of the given token
func (accounts *accounts) GetEthereumAccount() beth.Account {
	return accounts.ethereumAccount
}
