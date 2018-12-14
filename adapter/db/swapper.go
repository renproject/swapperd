package db

import (
	"encoding/base64"
	"encoding/json"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (db *dbStorage) PutReceipt(receipt swap.SwapReceipt) error {
	receiptData, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	id, err := base64.StdEncoding.DecodeString(string(receipt.ID))
	if err != nil {
		return err
	}
	return db.db.Put(append(TableSwapReceipts[:], id...), receiptData, nil)
}

func (db *dbStorage) UpdateReceipt(swapID swap.SwapID, update func(receipt *swap.SwapReceipt)) error {
	id, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return err
	}
	receiptBytes, err := db.db.Get(append(TableSwapReceipts[:], id...), nil)
	if err != nil {
		return err
	}
	receipt := swap.SwapReceipt{}
	if err := json.Unmarshal(receiptBytes, &receipt); err != nil {
		return err
	}
	update(&receipt)
	updatedReceiptBytes, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	return db.db.Put(append(TableSwapReceipts[:], id...), updatedReceiptBytes, nil)
}

func (db *dbStorage) Receipts() ([]swap.SwapReceipt, error) {
	iterator := db.db.NewIterator(&util.Range{Start: TableSwapReceiptsStart[:], Limit: TableSwapReceiptsLimit[:]}, nil)
	defer iterator.Release()
	receipts := []swap.SwapReceipt{}
	for iterator.Next() {
		value := iterator.Value()
		receipt := swap.SwapReceipt{}
		if err := json.Unmarshal(value, &receipt); err != nil {
			return receipts, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, iterator.Error()
}

func (db *dbStorage) Receipt(swapID swap.SwapID) (swap.SwapReceipt, error) {
	receipt := swap.SwapReceipt{}
	id, err := base64.StdEncoding.DecodeString(string(swapID))
	if err != nil {
		return receipt, err
	}

	receiptBytes, err := db.db.Get(append(TableSwapReceipts[:], id...), nil)
	if err != nil {
		return receipt, err
	}

	if err := json.Unmarshal(receiptBytes, &receipt); err != nil {
		return receipt, err
	}
	return receipt, nil
}

func (db *dbStorage) LoadCosts(swapID swap.SwapID) (blockchain.Cost, blockchain.Cost) {
	receipt, err := db.Receipt(swapID)
	if err != nil {
		return blockchain.Cost{}, blockchain.Cost{}
	}
	return receipt.SendCost, receipt.ReceiveCost
}
