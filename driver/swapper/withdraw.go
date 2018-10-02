package swapper

import (
	"errors"
	"fmt"
	"math"
	"math/big"
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
	// Parse and validate the token
	tokenStr = strings.ToLower(strings.TrimSpace(tokenStr))
	switch tokenStr {
	case "btc", "bitcoin", "xbt":
		valueBig, _ := big.NewFloat(value * math.Pow10(BTCDecimals)).Int(nil)
		feeBig, _ := big.NewFloat(fee * math.Pow10(BTCDecimals)).Int(nil)
		return swapper.withdrawBitcoin(to, valueBig.Int64(), feeBig.Int64())
	case "eth", "ethereum", "ether":
		valueBig, _ := big.NewFloat(value * math.Pow10(ETHDecimals)).Int(nil)
		feeBig, _ := big.NewFloat(fee * math.Pow10(BTCDecimals)).Int(nil)
		return swapper.withdrawEthereum(to, valueBig, feeBig)
	default:
		return errors.New("unknown token")
	}
}

func (swapper *swapper) withdrawBitcoin(to string, value, fee int64) error {
	conn := bitcoin.NewConnWithConfig(swapper.conf.Bitcoin)
	if fee == 0 {
		fee = 3000
	}
	return conn.Withdraw(to, swapper.keys.GetKey(token.BTC).(keystore.BitcoinKey), value, fee)
}

func (swapper *swapper) withdrawEthereum(to string, value, fee *big.Int) error {
	// Validate the receiver address
	if !strings.HasPrefix(to, "0x") {
		to = "0x" + to
	}
	if len(to) != 42 {
		return errors.New("invalid receiver address")
	}

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
