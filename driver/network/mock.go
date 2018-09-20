package network

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type mock struct {
	traderAddresses   map[order.ID][]byte
	traderSwapDetails map[order.ID][]byte
	mu                *sync.RWMutex
}

func NewMock() swap.Network {
	return &mock{
		traderAddresses:   map[order.ID][]byte{},
		traderSwapDetails: map[order.ID][]byte{},
		mu:                new(sync.RWMutex),
	}
}

func (mock *mock) SendOwnerAddress(orderID order.ID, address []byte) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()
	mock.traderAddresses[orderID] = address
	return nil
}

func (mock *mock) SendSwapDetails(orderID order.ID, swapDetails []byte) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()
	mock.traderSwapDetails[orderID] = swapDetails
	return nil
}

func (mock *mock) ReceiveOwnerAddress(orderID order.ID, waitTill int64) ([]byte, error) {
	for {
		mock.mu.Lock()
		details := mock.traderAddresses[orderID]
		mock.mu.Unlock()
		if bytes.Compare(details, []byte{}) != 0 {
			return details, nil
		}
		if time.Now().Unix() > waitTill {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return []byte{}, fmt.Errorf("Timed Out")
}

func (mock *mock) ReceiveSwapDetails(orderID order.ID, waitTill int64) ([]byte, error) {
	for {
		mock.mu.Lock()
		details := mock.traderSwapDetails[orderID]
		mock.mu.Unlock()
		if bytes.Compare(details, []byte{}) != 0 {
			return details, nil
		}
		if time.Now().Unix() > waitTill {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return []byte{}, fmt.Errorf("Timed Out")
}
