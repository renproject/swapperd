package state

import (
	"encoding/json"
	"sync"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type ActiveSwapList interface {
	AddSwap(orderID [32]byte) error
	DeleteSwap(orderID [32]byte) error
	PendingSwaps() [][32]byte
	ActiveSwaps() [][32]byte
	DeleteIfSettled(orderID [32]byte) error
	DeleteIfExpired(orderID [32]byte) error
}

type ProtectedSwapList struct {
	mu   *sync.RWMutex
	List [][32]byte `json:"list"`
}

func (list *ProtectedSwapList) Add(item [32]byte) {
	list.mu.Lock()
	defer list.mu.Unlock()
	list.List = append(list.List, item)
}

func (list *ProtectedSwapList) Delete(item [32]byte) {
	list.mu.Lock()
	defer list.mu.Unlock()
	for i, swap := range list.List {
		if swap == item {
			list.List = append(list.List[:i], list.List[i+1:]...)
			break
		}
	}
}

func NewProtectedSwapList() ProtectedSwapList {
	return ProtectedSwapList{
		mu:   new(sync.RWMutex),
		List: [][32]byte{},
	}
}

func (state *state) AddSwap(orderID [32]byte) error {
	defer state.LogInfo(orderID, "adding swap to the active list")
	pendingSwaps := NewProtectedSwapList()
	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err == nil {
		if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
			return err
		}
	}
	pendingSwaps.Add(orderID)
	if err := state.PutAddedAtTimestamp(orderID); err != nil {
		return err
	}
	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}
	return state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes)
}

func (state *state) DeleteSwap(orderID [32]byte) error {
	defer state.LogInfo(orderID, "removing swap from the active list")
	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return err
	}
	pendingSwaps := NewProtectedSwapList()
	if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
		return err
	}
	pendingSwaps.Delete(orderID)
	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}
	return state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes)
}

func (state *state) PendingSwaps() [][32]byte {
	pendingSwapsBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return nil
	}
	pendingSwaps := NewProtectedSwapList()
	if err := json.Unmarshal(pendingSwapsBytes, &pendingSwaps); err != nil {
		return nil
	}
	return pendingSwaps.List
}

func (state *state) ActiveSwaps() [][32]byte {
	exectableSwaps := [][32]byte{}
	for _, pendingSwap := range state.PendingSwaps() {
		if state.Status(pendingSwap) == swap.StatusExpired {
			continue
		}
		exectableSwaps = append(exectableSwaps, pendingSwap)
	}
	return exectableSwaps
}

func (state *state) DeleteIfExpired(orderID [32]byte) error {
	if state.Status(orderID) == swap.StatusExpired {
		return state.DeleteSwap(orderID)
	}
	return nil
}

func (state *state) DeleteIfSettled(orderID [32]byte) error {
	if state.Status(orderID) == swap.StatusSettled {
		return state.DeleteSwap(orderID)
	}
	return nil
}
