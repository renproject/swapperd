package state

import "encoding/json"

func (state *state) deleteSwap(orderID [32]byte) error {
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	defer state.LogInfo(orderID, "deleted the swap and its details")

	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return err
	}

	pendingSwaps := PendingSwaps{}
	if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
		return err
	}

	for i, swap := range pendingSwaps.Swaps {
		if swap == orderID {
			if len(pendingSwaps.Swaps) == 1 {
				pendingSwaps.Swaps = [][32]byte{}
				break
			}

			if i == 0 {
				pendingSwaps.Swaps = pendingSwaps.Swaps[1:]
				break
			}

			pendingSwaps.Swaps = append(pendingSwaps.Swaps[:i-1], pendingSwaps.Swaps[i:]...)
			break
		}
	}

	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}

	return state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes)
}

func (state *state) pendingSwaps() ([][32]byte, error) {
	pendingSwapsBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return [][32]byte{}, nil
	}

	pendingSwaps := PendingSwaps{}
	if err := json.Unmarshal(pendingSwapsBytes, &pendingSwaps); err != nil {
		return nil, err
	}

	return pendingSwaps.Swaps, nil
}
