package renex

import (
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type mockBinder struct {
	orderbook map[[32]byte]swap.Match
}

func NewMockBinder(matches ...swap.Match) (Binder, error) {
	orderbook := map[[32]byte]swap.Match{}
	for _, match := range matches {
		orderbook[match.PersonalOrderID] = match
	}
	return &mockBinder{
		orderbook: orderbook,
	}, nil
}

// GetOrderMatch gets the order match from the mock renex adapter.
func (binder *mockBinder) GetOrderMatch(orderID [32]byte, waitTill int64) (swap.Match, error) {
	return binder.orderbook[orderID], nil
}

// SubmitOrderMatch submits an order match to the mock binder adapter.
func (binder *mockBinder) SubmitOrderMatch(match swap.Match) error {
	binder.orderbook[match.PersonalOrderID] = match
	return nil
}
