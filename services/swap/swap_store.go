package swap

import (
	"github.com/republicprotocol/atom-go/services/store"
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
	return str.Write(append([]byte("status:"), orderID[:]...), []byte(status))
}

func (str *swapStore) ReadStatus(orderID [32]byte) (string, error) {
	status, err := str.Read(append([]byte("status:"), orderID[:]...))
	return string(status), err
}
