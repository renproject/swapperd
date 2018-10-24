package storage

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

var (
	TableSwaps      = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableSwapsStart = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableSwapsLimit = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	TablePendingSwaps      = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	TablePendingSwapsStart = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TablePendingSwapsLimit = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

type storage struct {
	mu *sync.RWMutex
	db *leveldb.DB
}

func New(db *leveldb.DB) swapper.Storage {
	return &storage{
		mu: new(sync.RWMutex),
		db: db,
	}
}

func (storage *storage) InsertSwap(swap swapper.Swap) error {
	pendingSwapData, err := json.Marshal(swap)
	if err != nil {
		return err
	}
	swapData, err := json.Marshal(swap.SwapBlob)
	if err != nil {
		return err
	}

	id, err := base64.StdEncoding.DecodeString(string(swap.ID))
	if err != nil {
		return err
	}

	if err := storage.db.Put(append(TablePendingSwaps[:], id...), pendingSwapData, nil); err != nil {
		return err
	}

	return storage.db.Put(append(TableSwaps[:], id...), swapData, nil)
}

func (storage *storage) DeletePendingSwap(swapID foundation.SwapID) error {
	id, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return err
	}
	return storage.db.Delete(append(TablePendingSwaps[:], id...), nil)
}

func (storage *storage) PendingSwap(swapID foundation.SwapID) (swapper.Swap, error) {
	id, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return swapper.Swap{}, err
	}
	swapBlobBytes, err := storage.db.Get(append(TablePendingSwaps[:], id...), nil)
	if err != nil {
		return swapper.Swap{}, err
	}
	swap := swapper.Swap{}
	if err := json.Unmarshal(swapBlobBytes, &swap); err != nil {
		return swapper.Swap{}, err
	}
	return swap, nil
}

func (storage *storage) Swaps() ([]foundation.SwapBlob, error) {
	iterator := storage.db.NewIterator(&util.Range{Start: TableSwapsStart[:], Limit: TableSwapsLimit[:]}, nil)
	defer iterator.Release()

	swaps := []foundation.SwapBlob{}
	for iterator.Next() {
		value := iterator.Value()
		swap := foundation.SwapBlob{}
		if err := json.Unmarshal(value, &swap); err != nil {
			return swaps, err
		}
		swaps = append(swaps, swap)
	}

	return swaps, iterator.Error()
}

func (storage *storage) PendingSwaps() ([]swapper.Swap, error) {
	iterator := storage.db.NewIterator(&util.Range{Start: TablePendingSwapsStart[:], Limit: TablePendingSwapsLimit[:]}, nil)
	defer iterator.Release()

	pendingSwaps := []swapper.Swap{}
	for iterator.Next() {
		value := iterator.Value()
		swap := swapper.Swap{}
		if err := json.Unmarshal(value, &swap); err != nil {
			return pendingSwaps, err
		}
		pendingSwaps = append(pendingSwaps, swap)
	}

	return pendingSwaps, iterator.Error()
}
