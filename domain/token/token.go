package token

import (
	"fmt"
)

var ErrUnsupportedToken = fmt.Errorf("Unsupported Token Code")

type Token string

const (
	BTC = Token("BTC")
	ETH = Token("ETH")
)

func TokenCodeToToken(tokenCode uint32) (Token, error) {
	switch tokenCode {
	case 0:
		return BTC, nil
	case 1:
		return ETH, nil
	default:
		return Token(""), ErrUnsupportedToken
	}
}
