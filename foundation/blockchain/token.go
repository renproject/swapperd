package blockchain

import (
	"fmt"
	"strings"
)

type ErrUnsupportedToken string

func NewErrUnsupportedToken(token TokenName) error {
	return ErrUnsupportedToken(fmt.Sprintf("unsupported token: %s", token))
}

func (err ErrUnsupportedToken) Error() string {
	return string(err)
}

type TokenName string

type Token struct {
	Name       TokenName      `json:"name"`
	Blockchain BlockchainName `json:"blockchain"`
}

func (token Token) String() string {
	return string(token.Name)
}

var (
	BTC  = TokenName("BTC")
	ETH  = TokenName("ETH")
	WBTC = TokenName("WBTC")
)

var (
	TokenBTC  = Token{TokenName("BTC"), Bitcoin}
	TokenETH  = Token{TokenName("ETH"), Ethereum}
	TokenWBTC = Token{TokenName("WBTC"), Ethereum}
)

func PatchToken(token string) (Token, error) {
	token = strings.ToLower(token)
	switch token {
	case "bitcoin", "btc", "xbt":
		return TokenBTC, nil
	case "wrappedbtc", "wbtc", "wrappedbitcoin":
		return TokenWBTC, nil
	case "ethereum", "eth", "ether":
		return TokenETH, nil
	default:
		return Token{}, fmt.Errorf("unsupported token: %s", token)
	}
}
