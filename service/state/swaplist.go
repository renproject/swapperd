package state

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type ActiveSwapList interface {
	AddSwap(orderID [32]byte) error
	DeleteSwap(orderID [32]byte) error
	PendingSwaps() ([][32]byte, error)
	ExecutableSwaps() ([][32]byte, error)
	ExpiredSwaps() ([][32]byte, error)
	DeleteIfRefunded(orderID [32]byte) error
	DeleteIfRedeemedOrExpired(orderID [32]byte) error
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
	if err := state.PutAddTimestamp(orderID); err != nil {
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

func (state *state) PendingSwaps() ([][32]byte, error) {
	pendingSwapsBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return [][32]byte{}, nil
	}
	pendingSwaps := NewProtectedSwapList()
	if err := json.Unmarshal(pendingSwapsBytes, &pendingSwaps); err != nil {
		return nil, err
	}
	return pendingSwaps.List, nil
}

func (state *state) ExecutableSwaps() ([][32]byte, error) {
	exectableSwaps := [][32]byte{}
	pendingSwaps, err := state.PendingSwaps()
	if err != nil {
		return nil, err
	}
	for _, pendingSwap := range pendingSwaps {
		if state.Status(pendingSwap) == swap.StatusExpired {
			continue
		}
		exectableSwaps = append(exectableSwaps, pendingSwap)
	}
	return exectableSwaps, nil
}

func (state *state) ExpiredSwaps() ([][32]byte, error) {
	pendingSwaps, err := state.PendingSwaps()
	if err != nil {
		return nil, err
	}
	refundableSwaps := [][32]byte{}
	for _, pendingSwap := range pendingSwaps {
		swapDet := state.ReadSwapDetails(pendingSwap)
		if swapDet.Request.TimeLock < time.Now().Unix() {
			refundableSwaps = append(refundableSwaps, pendingSwap)
		}
	}
	return refundableSwaps, nil
}

// TODO: check timestamp and delete if expired
func (state *state) DeleteIfRefunded(orderID [32]byte) error {
	if state.Status(orderID) == swap.StatusExpired {
		return state.DeleteSwap(orderID)
	}
	return nil
}

func (state *state) DeleteIfRedeemedOrExpired(orderID [32]byte) error {
	if state.Status(orderID) == swap.StatusExpired || state.Status(orderID) == swap.StatusSettled {
		return state.DeleteSwap(orderID)
	}
	return nil
}
