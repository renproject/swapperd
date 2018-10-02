package swapper

import (
	"fmt"
	"strconv"

	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/driver/bitcoin"
)

func (swapper *swapper) Withdraw(tok, to, value, fee string) error {
	switch token.Token(tok) {
	case token.BTC:
		return swapper.withdrawBitcoin(to, value, fee)
	case token.ETH:
		return swapper.withdrawEthereum(to, value, fee)
	default:
		return token.ErrUnsupportedToken
	}
}

func (swapper *swapper) withdrawBitcoin(to, valueString, feeString string) error {
	conn := bitcoin.NewConnWithConfig(swapper.conf.Bitcoin)
	val, err := strconv.ParseInt(valueString, 10, 64)
	if err != nil {
		val = 0
	}

	fee, err := strconv.ParseInt(feeString, 10, 64)
	if err != nil {
		fee = 3000
	}
	return conn.Withdraw(to, swapper.keys.GetKey(token.BTC).(keystore.BitcoinKey), val, fee)
}

func (swapper *swapper) withdrawEthereum(to, value, fee string) error {
	return fmt.Errorf("Ethereum Withdrawal Unimplemented")
	// conn, err := eth.NewConnWithConfig(swapper.Ethereum)
	// if err != nil {
	// 	return err
	// }

	// val, err := strconv.ParseInt(value, 10, 64)
	// if err != nil {
	// 	val = 0
	// }
	// return nil
}
