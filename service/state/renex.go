package state

import (
	"encoding/json"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type RenEx interface {
	AddSwap([32]byte) error
	ExecutableSwaps() ([][32]byte, error)
	DeleteIfRedeemedOrExpired([32]byte) error
}

func (state *state) AddSwap(orderID [32]byte) error {
	pendingSwaps := PendingSwaps{}
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err == nil {
		if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
			return err
		}
	}
	pendingSwaps.Swaps = append(pendingSwaps.Swaps, orderID)
	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}
	if err := state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes); err != nil {
		return err
	}
	return state.PutOrderTimeStamp(orderID)
}

func (state *state) ExecutableSwaps() ([][32]byte, error) {
	exectableSwaps := [][32]byte{}
	state.swapMu.RLock()
	pendingSwaps, err := state.pendingSwaps()
	state.swapMu.RUnlock()
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

func (state *state) DeleteIfRedeemedOrExpired(orderID [32]byte) error {
	if state.Status(orderID) != swap.StatusRedeemed && state.Status(orderID) != swap.StatusExpired {
		return nil
	}
	return state.deleteSwap(orderID)
}
