package match

import (
	"math/big"
)

type Match interface {
	PersonalOrderID() [32]byte
	ForeignOrderID() [32]byte
	SendValue() *big.Int
	RecieveValue() *big.Int
	SendCurrency() uint32
	RecieveCurrency() uint32
}

type match struct {
	personalOrderID [32]byte
	foreignOrderID  [32]byte
	sendValue       *big.Int
	recieveValue    *big.Int
	sendCurrency    uint32
	recieveCurrency uint32
}

func NewMatch(personalOrderID, foreignOrderID [32]byte, sendValue, recieveValue *big.Int, sendCurrency, recieveCurrency uint32) Match {
	return &match{
		personalOrderID: personalOrderID,
		foreignOrderID:  foreignOrderID,
		sendValue:       sendValue,
		recieveValue:    recieveValue,
		sendCurrency:    sendCurrency,
		recieveCurrency: recieveCurrency,
	}
}

func (match *match) PersonalOrderID() [32]byte {
	return match.personalOrderID
}
func (match *match) ForeignOrderID() [32]byte {
	return match.foreignOrderID
}
func (match *match) SendValue() *big.Int {
	return match.sendValue
}
func (match *match) RecieveValue() *big.Int {
	return match.recieveValue
}
func (match *match) SendCurrency() uint32 {
	return match.sendCurrency
}
func (match *match) RecieveCurrency() uint32 {
	return match.recieveCurrency
}
