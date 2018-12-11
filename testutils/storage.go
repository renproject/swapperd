package testutils

//
// import (
// 	"errors"
// 	"sync"
//
// 	"github.com/republicprotocol/swapperd/foundation/swap"
// )
//
// type MockStorage struct {
// 	mu     *sync.RWMutex
// 	swaps  map[swap.SwapID]swap.SwapRequest
// 	status map[swap.SwapID]swap.SwapStatus
// }
//
// func NewMockStorage() MockStorage {
// 	return MockStorage{
// 		mu:    new(sync.RWMutex),
// 		swaps: map[swap.SwapID]swap.SwapRequest{},
// 	}
// }
//
// func (store *MockStorage) InsertSwap(swap swap.SwapRequest) error {
// 	store.mu.Lock()
// 	defer store.mu.Unlock()
//
// 	store.swaps[swap.ID] = swap
// 	store.status[swap.ID] = swap.NewSwapStatus(swap.SwapBlob)
// 	return nil
// }
//
// func (store *MockStorage) PendingSwap(id swap.SwapID) (swap.SwapRequest, error) {
// 	store.mu.RLock()
// 	defer store.mu.RUnlock()
//
// 	swap, ok := store.swaps[id]
// 	if !ok {
// 		return swap.SwapRequest{}, errors.New("pending swap not found")
// 	}
// 	return swap, nil
// }
//
// func (store *MockStorage) DeletePendingSwap(id swap.SwapID) error {
// 	store.mu.Lock()
// 	defer store.mu.Unlock()
//
// 	delete(store.swaps, id)
// 	return nil
// }
//
// func (store *MockStorage) PendingSwaps() ([]swap.SwapRequest, error) {
// 	store.mu.Lock()
// 	defer store.mu.Unlock()
//
// 	swaps := make([]swap.SwapRequest, 0, len(store.status))
// 	for _, status := range store.status {
// 		if status.Status == swap.Initiated {
// 			swaps = append(swaps, store.swaps[status.ID])
// 		}
// 	}
//
// 	return swaps, nil
// }
//
// func (store *MockStorage) UpdateStatus(update swap.StatusUpdate) error {
// 	store.mu.Lock()
// 	defer store.mu.Unlock()
//
// 	status, ok := store.status[update.ID]
// 	if !ok {
// 		return errors.New("swap not found")
// 	}
// 	status.Status = update.Status
// 	store.status[update.ID] = status
// 	return nil
// }
//
// func (store *MockStorage) Swaps() ([]swap.SwapStatus, error) {
// 	store.mu.RLock()
// 	defer store.mu.RUnlock()
//
// 	statuses := make([]swap.SwapStatus, 0, len(store.status))
// 	for _, status := range store.status {
// 		statuses = append(statuses, status)
// 	}
//
// 	return statuses, nil
// }
