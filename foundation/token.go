package foundation

import "fmt"

func NewErrUnsupportedToken(token string) error {
	return fmt.Errorf("unsupported token: %s", token)
}

type Token struct {
	Name       string
	Blockchain string
}

func (token Token) String() string {
	return token.Name
}

var (
	TokenBTC  = Token{"BTC", "bitcoin"}
	TokenETH  = Token{"ETH", "ethereum"}
	TokenWBTC = Token{"WBTC", "ethereum"}
)
