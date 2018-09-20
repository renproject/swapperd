package state

import (
	"encoding/json"
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

type SwapList struct {
	List [][32]byte `json:"list"`
}

func (list *SwapList) Add(item [32]byte) {
	list.List = append(list.List, item)
}

func (list *SwapList) Delete(item [32]byte) {
	for i, swap := range list.List {
		if swap == item {
			list.List = append(list.List[:i], list.List[i+1:]...)
			break
		}
	}
}

func (state *state) AddSwap(orderID [32]byte) error {
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	defer state.LogInfo(orderID, "adding swap to the active list")
	pendingSwaps := SwapList{}
	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err == nil {
		if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
			return err
		}
	}
	pendingSwaps.Add(orderID)
	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}
	return state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes)
}

func (state *state) DeleteSwap(orderID [32]byte) error {
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	defer state.LogInfo(orderID, "removing swap from the active list")
	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return err
	}
	pendingSwaps := SwapList{}
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
	state.swapMu.RLock()
	defer state.swapMu.RUnlock()
	pendingSwapsBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return [][32]byte{}, nil
	}
	pendingSwaps := SwapList{}
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
		if state.Status(pendingSwap) == swap.StatusComplained {
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
		_, expiry, err := state.InitiateDetails(pendingSwap)
		if err != nil {
			continue
		}
		if expiry < time.Now().Unix() {
			refundableSwaps = append(refundableSwaps, pendingSwap)
		}
	}
	return refundableSwaps, nil
}

func (state *state) DeleteIfRefunded(orderID [32]byte) error {
	if state.Status(orderID) == swap.StatusRefunded {
		return state.DeleteSwap(orderID)
	}
	return nil
}

func (state *state) DeleteIfRedeemedOrExpired(orderID [32]byte) error {
	if state.Status(orderID) == swap.StatusRedeemed || state.Status(orderID) == swap.StatusExpired {
		return state.DeleteSwap(orderID)
	}
	return nil
}
