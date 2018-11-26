package db

import (
	"encoding/base64"
	"encoding/json"

	"github.com/republicprotocol/swapperd/adapter/router"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	TableSwaps      = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableSwapsStart = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableSwapsLimit = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	TablePendingSwaps      = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	TablePendingSwapsStart = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TablePendingSwapsLimit = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

type dbStorage struct {
	db *leveldb.DB
}

func New(db *leveldb.DB) router.Storage {
	return &dbStorage{
		db: db,
	}
}

func (db *dbStorage) InsertSwap(swap foundation.SwapRequest) error {
	pendingSwapData, err := json.Marshal(swap)
	if err != nil {
		return err
	}
	swapData, err := json.Marshal(foundation.NewSwapStatus(swap.SwapBlob))
	if err != nil {
		return err
	}
	id, err := base64.StdEncoding.DecodeString(string(swap.ID))
	if err != nil {
		return err
	}
	if err := db.db.Put(append(TablePendingSwaps[:], id...), pendingSwapData, nil); err != nil {
		return err
	}
	return db.db.Put(append(TableSwaps[:], id...), swapData, nil)
}

func (db *dbStorage) DeletePendingSwap(swapID foundation.SwapID) error {
	id, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return err
	}
	return db.db.Delete(append(TablePendingSwaps[:], id...), nil)
}

func (db *dbStorage) PendingSwap(swapID foundation.SwapID) (foundation.SwapRequest, error) {
	id, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return foundation.SwapRequest{}, err
	}
	swapBlobBytes, err := db.db.Get(append(TablePendingSwaps[:], id...), nil)
	if err != nil {
		return foundation.SwapRequest{}, err
	}
	swap := foundation.SwapRequest{}
	if err := json.Unmarshal(swapBlobBytes, &swap); err != nil {
		return foundation.SwapRequest{}, err
	}
	return swap, nil
}

func (db *dbStorage) Swaps() ([]foundation.SwapStatus, error) {
	iterator := db.db.NewIterator(&util.Range{Start: TableSwapsStart[:], Limit: TableSwapsLimit[:]}, nil)
	defer iterator.Release()
	swaps := []foundation.SwapStatus{}
	for iterator.Next() {
		value := iterator.Value()
		swap := foundation.SwapStatus{}
		if err := json.Unmarshal(value, &swap); err != nil {
			return swaps, err
		}
		swaps = append(swaps, swap)
	}
	return swaps, iterator.Error()
}

func (db *dbStorage) PendingSwaps() ([]foundation.SwapRequest, error) {
	iterator := db.db.NewIterator(&util.Range{Start: TablePendingSwapsStart[:], Limit: TablePendingSwapsLimit[:]}, nil)
	defer iterator.Release()
	pendingSwaps := []foundation.SwapRequest{}
	for iterator.Next() {
		value := iterator.Value()
		swap := foundation.SwapRequest{}
		if err := json.Unmarshal(value, &swap); err != nil {
			return pendingSwaps, err
		}
		pendingSwaps = append(pendingSwaps, swap)
	}
	return pendingSwaps, iterator.Error()
}

func (db *dbStorage) UpdateStatus(update foundation.StatusUpdate) error {
	id, err := base64.StdEncoding.DecodeString(string(update.ID))
	if err != nil {
		return err
	}
	receiptBytes, err := db.db.Get(append(TableSwaps[:], id...), nil)
	if err != nil {
		return err
	}
	status := foundation.SwapStatus{}
	if err := json.Unmarshal(receiptBytes, &status); err != nil {
		return err
	}
	status.Status = update.Status
	updatedReceiptBytes, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return db.db.Put(append(TableSwaps[:], id...), updatedReceiptBytes, nil)
}