package fund

import (
	"math/big"

	beth "github.com/republicprotocol/beth-go"
	libbtc "github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/foundation"
)

type Balance struct {
	Address string
	Amount  *big.Int
}

type Manager interface {
	SupportedTokens() ([]foundation.Token, error)
	Balances() (map[foundation.Token]Balance, error)
	Withdraw(password string, token foundation.Token, to string, amount *big.Int) (string, error)
	GetBitcoinAccount(password string) (libbtc.Account, error)
	GetEthereumAccount(password string) (beth.Account, error)
}

type manager struct {
	supportedTokens []foundation.Token
	addresses       map[foundation.Blockchain]string
	ethAccount      beth.Account
	btcAccount      libbtc.Account
}

func New(config Config) (Manager, error) {
	supportedTokens, err := decodeSupportedTokens(config)
	if err != nil {
		return nil, err
	}
	addresses, err := decodeAddresses(config)
	if err != nil {
		return nil, err
	}
	ethAccount, btcAccount, err := generateRandomAccounts(config)
	if err != nil {
		return nil, err
	}
	return &manager{
		supportedTokens: supportedTokens,
		addresses:       addresses,
		ethAccount:      ethAccount,
		btcAccount:      btcAccount,
	}, nil
}

func generateRandomAccounts(config Config) (beth.Account, libbtc.Account, error) {
	return beth.NewAccount()
}

func decodeAddresses(config)
