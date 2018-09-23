package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/renex-swapper-go/service/renex"
)

type mock struct {
	traderSwapDetails map[[32]byte]renex.SwapDetails
	mu                *sync.RWMutex
}

func NewMock() renex.Network {
	return &mock{
		traderSwapDetails: map[[32]byte]renex.SwapDetails{},
		mu:                new(sync.RWMutex),
	}
}

func (mock *mock) SendSwapDetails(orderID [32]byte, swapDetails renex.SwapDetails) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()
	mock.traderSwapDetails[orderID] = swapDetails
	return nil
}

func (mock *mock) ReceiveSwapDetails(orderID [32]byte, waitTill int64) (renex.SwapDetails, error) {
	det := renex.SwapDetails{}
	for {
		mock.mu.RLock()
		details := mock.traderSwapDetails[orderID]
		mock.mu.RUnlock()
		if details != det {
			return details, nil
		}
		if time.Now().Unix() > waitTill {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return det, fmt.Errorf("Timed Out")
}
