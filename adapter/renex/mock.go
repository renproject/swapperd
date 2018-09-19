package renex

import (
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
)

type mockBinder struct {
	orderbook map[order.ID]match.Match
}

func NewMockBinder(matches ...match.Match) (Binder, error) {
	orderbook := map[order.ID]match.Match{}
	for _, match := range matches {
		orderbook[match.PersonalOrderID()] = match
	}
	return &mockBinder{
		orderbook: orderbook,
	}, nil
}

// GetOrderMatch gets the order match from the mock renex adapter.
func (binder *mockBinder) GetOrderMatch(orderID order.ID, waitTill int64) (match.Match, error) {
	return binder.orderbook[orderID], nil
}

// SubmitOrderMatch submits an order match to the mock binder adapter.
func (binder *mockBinder) SubmitOrderMatch(match match.Match) error {
	binder.orderbook[match.PersonalOrderID()] = match
	return nil
}
