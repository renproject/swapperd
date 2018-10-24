package foundation

import "fmt"

type ErrUnsupportedToken string

func NewErrUnsupportedToken(token string) error {
	return ErrUnsupportedToken(fmt.Sprintf("unsupported token: %s", token))
}

func (err ErrUnsupportedToken) Error() string {
	return string(err)
}

var (
	Bitcoin  = "bitcoin"
	Ethereum = "ethereum"
)

type Token struct {
	Name       string `json:"name"`
	Blockchain string `json:"blockchain"`
}

func (token Token) String() string {
	return token.Name
}

var (
	TokenBTC  = Token{"BTC", Bitcoin}
	TokenETH  = Token{"ETH", Ethereum}
	TokenWBTC = Token{"WBTC", Ethereum}
)
