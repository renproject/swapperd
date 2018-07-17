package store

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/republicprotocol/atom-go/domains/match"
)

type State interface {
	GetSwaps() [][32]byte
	AddSwap([32]byte) error
	DeleteSwap([32]byte) error

	UpdateStatus([32]byte, string) error
	ReadStatus([32]byte) string

	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error

	SetMatch([32]byte, match.Match) error
	GetMatch([32]byte) (match.Match, error)
}

type state struct {
	Store
	mu *sync.RWMutex
}

// PendingSwaps stores all the swaps that are pending
type PendingSwaps struct {
	Swaps [][32]byte `json:"pendingSwaps"`
}

func NewSwapStore(store Store) State {
	return &state{
		Store: store,
		mu:    new(sync.RWMutex),
	}
}

func (str *state) UpdateStatus(orderID [32]byte, status string) error {
	return str.Write(append([]byte("status:"), orderID[:]...), []byte(status))
}

func (str *state) ReadStatus(orderID [32]byte) string {
	status, err := str.Read(append([]byte("status:"), orderID[:]...))
	if err != nil {
		return "UNKNOWN"
	}
	return string(status)
}

func (str *state) SetMatch(orderID [32]byte, m match.Match) error {
	data, err := m.Serialize()
	if err != nil {
		return err
	}
	return str.Write(append([]byte("match:"), orderID[:]...), data)
}

func (str *state) GetMatch(orderID [32]byte) (match.Match, error) {
	data, err := str.Read(append([]byte("match:"), orderID[:]...))
	if err != nil {
		return nil, err
	}
	return match.NewMatchFromBytes(data)
}

func (str *state) AddSwap(swap [32]byte) error {
	swaps := str.GetSwaps()
	str.mu.Lock()
	defer str.mu.Unlock()
	swaps = append(swaps, swap)
	pending := PendingSwaps{
		Swaps: swaps,
	}
	swapData, err := json.Marshal(pending)
	if err != nil {
		return err
	}
	return str.Write(append([]byte("pending")), swapData)
}

func (str *state) DeleteSwap(swap [32]byte) error {
	swaps := str.GetSwaps()
	str.mu.Lock()
	defer str.mu.Unlock()
	for i, swapElement := range swaps {
		if swap == swapElement {
			swaps = append(swaps[:i], swaps[i+1:]...)
			swapData, err := json.Marshal(swaps)
			if err != nil {
				return err
			}
			return str.Write(append([]byte("pending")), swapData)
		}
	}
	return errors.New("Swap Not found")
}

func (str *state) GetSwaps() [][32]byte {
	str.mu.RLock()
	defer str.mu.RUnlock()
	var pending PendingSwaps
	pendingSwaps, err := str.Read(append([]byte("pending")))
	if err != nil {
		return [][32]byte{}
	}
	err = json.Unmarshal(pendingSwaps, &pending)
	if err != nil {
		return [][32]byte{}
	}
	return pending.Swaps
}
