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
	REN  = TokenName("REN")
	TUSD = TokenName("TUSD")
	DGX  = TokenName("DGX")
	ZRX  = TokenName("ZRX")
	OMG  = TokenName("OMG")
)

var (
	TokenBTC  = Token{TokenName("BTC"), Bitcoin}
	TokenETH  = Token{TokenName("ETH"), Ethereum}
	TokenWBTC = Token{TokenName("WBTC"), Ethereum}
	TokenREN  = Token{TokenName("REN"), Ethereum}
	TokenTUSD = Token{TokenName("TUSD"), Ethereum}
	TokenDGX  = Token{TokenName("DGX"), Ethereum}
	TokenZRX  = Token{TokenName("ZRX"), Ethereum}
	TokenOMG  = Token{TokenName("OMG"), Ethereum}
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
	case "ren", "republictoken", "republic token":
		return TokenREN, nil
	case "tusd", "trueusd", "true usd":
		return TokenTUSD, nil
	case "digix gold token", "dgx", "dgt":
		return TokenDGX, nil
	case "zerox", "zrx", "0x":
		return TokenZRX, nil
	case "omisego", "omg", "omise go":
		return TokenOMG, nil
	default:
		return Token{}, fmt.Errorf("unsupported token: %s", token)
	}
}

func IsValidToken(name TokenName) bool {
	return (name == BTC || name == ETH || name == WBTC)
}
