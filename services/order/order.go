package order

import "math/big"

type Order interface {
	PersonalOrderID() [32]byte
	ForeignOrderID() [32]byte
	SendValue() *big.Int
	RecieveValue() *big.Int
	SendCurrency() string
	RecieveCurrency() string
}
