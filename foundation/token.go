package foundation

import (
	"fmt"
	"strings"
)

// ErrUnsupportedToken is returned when the token is not supported by swapperd.
type ErrUnsupportedToken error

// NewErrUnsupportedToken returns a new ErrUnsupportedToken error.
func NewErrUnsupportedToken(token TokenName) error {
	return fmt.Errorf("unsupported token: %s", token)
}

// BlockchainName is the name of the blockchain.
type BlockchainName string

var (
	Bitcoin  = BlockchainName("bitcoin")
	Ethereum = BlockchainName("ethereum")
)

// TokenName is the name of the token.
type TokenName string

var (
	BTC  = TokenName("BTC")
	ETH  = TokenName("ETH")
	WBTC = TokenName("WBTC")
)

// Token represents the token we are trading.
type Token struct {
	Name       TokenName      `json:"name"`
	Blockchain BlockchainName `json:"blockchain"`
}

var (
	TokenBTC  = Token{TokenName("BTC"), Bitcoin}
	TokenETH  = Token{TokenName("ETH"), Ethereum}
	TokenWBTC = Token{TokenName("WBTC"), Ethereum}
)

func (token Token) String() string {
	return string(token.Name)
}

type Balance struct {
	Address string `json:"address"`
	Amount  string `json:"balance"`
}

type Blockchain struct {
	Name    BlockchainName `json:"name"`
	Address string         `json:"address"`
}

func PatchToken(token string) (Token, error) {
	token = strings.TrimSpace(strings.ToLower(token))
	switch token {
	case "bitcoin", "btc", "xbt":
		return TokenBTC, nil
	case "wrappedbtc", "wbtc", "wrappedbitcoin":
		return TokenWBTC, nil
	case "ethereum", "eth", "ether":
		return TokenETH, nil
	default:
		return Token{}, NewErrUnsupportedToken(TokenName(token))
	}
}