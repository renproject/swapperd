package wallet

import (
	"github.com/renproject/tokens"
)

func (wallet *wallet) SupportedTokens() []tokens.Token {
	return tokens.SupportedTokens
}
