package swapper

import (
	"math/big"

	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

func (swapper *Swapper) Withdraw(to string, token token.Token, value *big.Int) {
	switch token {
	case token.ETH:
		swapper.withdrawBitcoin(to string, value.Int64())
	case token.BTC:
	}
}
