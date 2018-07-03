package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/atom-go/services/swap"
)

type MockNetwork struct {
	mu          *sync.RWMutex
	SwapDetails map[[32]byte]chan []byte
}

func NewMockNetwork() swap.Network {
	return &MockNetwork{
		mu:          new(sync.RWMutex),
		SwapDetails: map[[32]byte]chan []byte{},
	}
}

func (net *MockNetwork) SendSwapDetails(orderID [32]byte, swapDetails []byte) error {
	swaps := make(chan []byte, 2)
	fmt.Println("Writing to channel")
	swaps <- swapDetails

	fmt.Println("Writing to map")
	net.mu.Lock()
	net.SwapDetails[orderID] = swaps
	net.mu.Unlock()
	return nil
}

func (net *MockNetwork) RecieveSwapDetails(orderID [32]byte) ([]byte, error) {
	fmt.Println("Reading from channel")

	for {
		time.Sleep(time.Second)
		net.mu.RLock()
		if _, ok := net.SwapDetails[orderID]; ok {
			net.mu.RUnlock()
			break
		}
		net.mu.RUnlock()
	}

	swapDetails := <-net.SwapDetails[orderID]
	fmt.Println("Read from channel")
	return swapDetails, nil
}
