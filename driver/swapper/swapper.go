package swapper

import (
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

type Swapper interface {
	Http(port int64)
	Withdraw(token.Token, string, string, string)
}
type swapper struct {
}

func NewSwapper() Swapper {
	return &swapper{}
}
