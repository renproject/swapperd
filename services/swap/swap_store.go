package swap

import (
	"fmt"

	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/republic-go/order"
)

type SwapStore interface {
	UpdateStatus([32]byte, string) error
	ReadStatus([32]byte) (string, error)
}

type swapStore struct {
	store.Store
}

func NewSwapStore(store store.Store) SwapStore {
	return &swapStore{
		store,
	}
}

func (str *swapStore) UpdateStatus(orderID [32]byte, status string) error {
	fmt.Printf("Order %v updated to status %v\n", order.ID(orderID), status)
	return str.Write(append([]byte("status:"), orderID[:]...), []byte(status))
}

func (str *swapStore) ReadStatus(orderID [32]byte) (string, error) {
	status, err := str.Read(append([]byte("status:"), orderID[:]...))
	return string(status), err
}
