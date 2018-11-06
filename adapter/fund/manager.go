package fund

import (
	"math/big"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/foundation"
)

type Balance struct {
	Address string
	Amount  *big.Int
}

type Manager interface {
	SupportedTokens() []foundation.Token
	SupportedBlockchains() []Blockchain
	Balances() (map[foundation.Token]Balance, error)
	Withdraw(password string, token foundation.Token, to string, amount *big.Int) (string, error)
	GetBitcoinAccount(password string) (libbtc.Account, error)
	GetEthereumAccount(password string) (beth.Account, error)
}

type manager struct {
	config Config
}

func New(config Config) Manager {
	return &manager{
		config: config,
	}
}
