package foundation

import (
	"fmt"
	"strings"
)

type ErrUnsupportedToken string

func NewErrUnsupportedToken(token string) error {
	return ErrUnsupportedToken(fmt.Sprintf("unsupported token: %s", token))
}

func (err ErrUnsupportedToken) Error() string {
	return string(err)
}

type Blockchain string

var (
	Bitcoin  = Blockchain("bitcoin")
	Ethereum = Blockchain("ethereum")
)

type Token struct {
	Name       string     `json:"name"`
	Blockchain Blockchain `json:"blockchain"`
}

func (token Token) String() string {
	return token.Name
}

var (
	TokenBTC  = Token{"BTC", Bitcoin}
	TokenETH  = Token{"ETH", Ethereum}
	TokenWBTC = Token{"WBTC", Ethereum}
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
