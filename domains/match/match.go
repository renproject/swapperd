package match

import (
	"math/big"
)

type Match interface {
	PersonalOrderID() [32]byte
	ForeignOrderID() [32]byte
	SendValue() *big.Int
	ReceiveValue() *big.Int
	SendCurrency() uint32
	ReceiveCurrency() uint32
}

type match struct {
	personalOrderID [32]byte
	foreignOrderID  [32]byte
	sendValue       *big.Int
	receiveValue    *big.Int
	sendCurrency    uint32
	receiveCurrency uint32
}

func NewMatch(personalOrderID, foreignOrderID [32]byte, sendValue, receiveValue *big.Int, sendCurrency, receiveCurrency uint32) Match {
	return &match{
		personalOrderID: personalOrderID,
		foreignOrderID:  foreignOrderID,
		sendValue:       sendValue,
		receiveValue:    receiveValue,
		sendCurrency:    sendCurrency,
		receiveCurrency: receiveCurrency,
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
func (match *match) ReceiveValue() *big.Int {
	return match.receiveValue
}
func (match *match) SendCurrency() uint32 {
	return match.sendCurrency
}
func (match *match) ReceiveCurrency() uint32 {
	return match.receiveCurrency
}
