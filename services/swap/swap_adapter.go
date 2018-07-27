package swap

import "github.com/republicprotocol/atom-go/domains/order"

type SwapAdapter interface {
	SendOwnerAddress(order.ID, []byte) error
	ReceiveOwnerAddress(order.ID) ([]byte, error)
	ReceiveSwapDetails(order.ID, bool) ([]byte, error)
	SendSwapDetails(order.ID, []byte) error
}
