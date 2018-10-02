package swapper

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/driver/bitcoin"
)

const (
	BTCDecimals = 8
	ETHDecimals = 18
)

func (swapper *swapper) Withdraw(tokenStr, to string, value, fee float64) error {
	// Validate the receiver address
	if !strings.HasPrefix(to, "0x") {
		to = "0x" + to
	}
	if len(to) != 42 {
		return errors.New("invalid receiver address")
	}

	// Parse and validate the token
	tokenStr = strings.ToLower(strings.TrimSpace(tokenStr))
	switch tokenStr {
	case "btc", "bitcoin", "xbt":
		valueBig, _ := big.NewFloat(value * math.Pow10(BTCDecimals)).Int(nil)
		feeBig, _ := big.NewFloat(fee * math.Pow10(BTCDecimals)).Int(nil)
		return swapper.withdrawBitcoin(to, valueBig, feeBig)
	case "eth", "ethereum", "ether":
		valueBig, _ := big.NewFloat(value * math.Pow10(ETHDecimals)).Int(nil)
		feeBig, _ := big.NewFloat(fee * math.Pow10(BTCDecimals)).Int(nil)
		return swapper.withdrawEthereum(to, valueBig, feeBig)
	default:
		return errors.New("unknown token")
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
