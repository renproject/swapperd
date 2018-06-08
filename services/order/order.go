package order

import "math/big"

type Order interface {
	MyOrderID() [32]byte
	TradingOrderID() [32]byte
	SendValue() *big.Int
	RecieveValue() *big.Int
	SendCurrency() string
	RecieveCurrency() string
}
