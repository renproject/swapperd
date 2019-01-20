package blockchain

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
)

// ErrUnsupportedToken is returned when the token is not supported by swapperd.
type ErrUnsupportedToken error

// NewErrUnsupportedToken returns a new ErrUnsupportedToken error.
func NewErrUnsupportedToken(token TokenName) error {
	return fmt.Errorf("unsupported token: %s", token)
}

type TokenName string

var (
	BTC  = TokenName("BTC")
	ETH  = TokenName("ETH")
	WBTC = TokenName("WBTC")
	REN  = TokenName("REN")
	DGX  = TokenName("DGX")
	ZRX  = TokenName("ZRX")
	OMG  = TokenName("OMG")
	DAI  = TokenName("DAI")
	PAX  = TokenName("PAX")
	USDC = TokenName("USDC")
	GUSD = TokenName("GUSD")
	TUSD = TokenName("TUSD")
)

// Token represents the token we are trading.
type Token struct {
	Name       TokenName      `json:"name"`
	Decimals   int            `json:"decimals"`
	Blockchain BlockchainName `json:"blockchain"`
}

var (
	TokenBTC  = Token{BTC, 8, Bitcoin}
	TokenETH  = Token{ETH, 18, Ethereum}
	TokenWBTC = Token{WBTC, 8, ERC20}
	TokenREN  = Token{REN, 18, ERC20}
	TokenDGX  = Token{DGX, 9, ERC20}
	TokenZRX  = Token{ZRX, 18, ERC20}
	TokenOMG  = Token{OMG, 18, ERC20}
	TokenTUSD = Token{TUSD, 18, ERC20}
	TokenDAI  = Token{DAI, 18, ERC20}
	TokenPAX  = Token{PAX, 18, ERC20}
	TokenUSDC = Token{USDC, 6, ERC20}
	TokenGUSD = Token{GUSD, 2, ERC20}
)

var SupportedTokens = []Token{
	TokenBTC, TokenETH, TokenWBTC, TokenREN, TokenDGX, TokenZRX, TokenOMG, TokenTUSD, TokenDAI, TokenUSDC, TokenGUSD, TokenPAX,
}

func (token Token) String() string {
	return string(token.Name)
}

func (Token) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(SupportedTokens[rand.Int()%len(SupportedTokens)])
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
	case "ren", "republictoken", "republic token":
		return TokenREN, nil
	case "tusd", "trueusd", "true-usd":
		return TokenTUSD, nil
	case "digix-gold-token", "dgx", "dgt":
		return TokenDGX, nil
	case "zerox", "zrx", "0x":
		return TokenZRX, nil
	case "omisego", "omg", "omise-go":
		return TokenOMG, nil
	case "usdc", "usd-coin", "usdcoin":
		return TokenUSDC, nil
	case "dai", "maker-dai", "makerdai":
		return TokenDAI, nil
	case "gusd", "gemini-dollar", "geminidollar":
		return TokenGUSD, nil
	case "pax", "paxos-standard-token", "paxosstandardtoken":
		return TokenPAX, nil
	default:
		return Token{}, NewErrUnsupportedToken(TokenName(token))
	}
}
