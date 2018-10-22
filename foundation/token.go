package foundation

import "fmt"

type ErrUnsupportedToken string

func NewErrUnsupportedToken(token string) error {
	return ErrUnsupportedToken(fmt.Sprintf("unsupported token: %s", token))
}

func (err ErrUnsupportedToken) Error() string {
	return string(err)
}

type Token struct {
	Name       string `json:"name"`
	Blockchain string `json:"blockchain"`
}

func (token Token) String() string {
	return token.Name
}

var (
	TokenBTC  = Token{"BTC", "bitcoin"}
	TokenETH  = Token{"ETH", "ethereum"}
	TokenWBTC = Token{"WBTC", "ethereum"}
)
