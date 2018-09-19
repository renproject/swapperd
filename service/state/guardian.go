package state

import (
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type Guardian interface {
	RefundableSwaps() ([][32]byte, error)
	DeleteIfRefunded([32]byte) error
}

func (state *state) RefundableSwaps() ([][32]byte, error) {
	state.swapMu.RLock()
	pendingSwaps, err := state.pendingSwaps()
	state.swapMu.RUnlock()
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
	if state.Status(orderID) != swap.StatusRefunded {
		return nil
	}
	return state.deleteSwap(orderID)
}
