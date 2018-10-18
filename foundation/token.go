package foundation

import "fmt"

func NewErrUnsupportedToken(token string) error {
	return fmt.Errorf("unsupported token: %s", token)
}

type Token struct {
	Name       string
	Blockchain string
}

var (
	TokenBTC  = Token{"BTC", "Bitcoin"}
	TokenETH  = Token{"ETH", "Ethereum"}
	TokenWBTC = Token{"WBTC", "Ethereum"}
)
