package mock

import (
	"math/big"

	"github.com/republicprotocol/atom-go/services/order"
)

type MockOrder struct {
	personalOrderID [32]byte
	foreignOrderID  [32]byte
	sendValue       *big.Int
	recieveValue    *big.Int
	sendCurrency    string
	recieveCurrency string
}

func NewMockOrder(personalOrderID, foreignOrderID [32]byte, sendValue, recieveValue *big.Int, sendCurrency, recieveCurrency string) order.Order {
	return &MockOrder{
		personalOrderID: personalOrderID,
		foreignOrderID:  foreignOrderID,
		sendValue:       sendValue,
		recieveValue:    recieveValue,
		sendCurrency:    sendCurrency,
		recieveCurrency: recieveCurrency,
	}
}

func (order *MockOrder) PersonalOrderID() [32]byte {
	return order.personalOrderID
}
func (order *MockOrder) ForeignOrderID() [32]byte {
	return order.foreignOrderID
}
func (order *MockOrder) SendValue() *big.Int {
	return order.sendValue
}
func (order *MockOrder) RecieveValue() *big.Int {
	return order.recieveValue
}
func (order *MockOrder) SendCurrency() string {
	return order.sendCurrency
}
func (order *MockOrder) RecieveCurrency() string {
	return order.recieveCurrency
}
