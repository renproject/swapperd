package db

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/republicprotocol/swapperd/core/transfer"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	TableTransfer      = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}
	TableTransferStart = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableTransferLimit = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

func (db *dbStorage) PutTransfer(transfer transfer.TransferReceipt) error {
	transferData, err := json.Marshal(transfer)
	if err != nil {
		return err
	}
	txHashBytes, err := txHashToBytes(transfer.TxHash)
	if err != nil {
		return err
	}
	return db.db.Put(append(TableTransfer[:], txHashBytes...), transferData, nil)
}

func (db *dbStorage) Transfers() ([]transfer.TransferReceipt, error) {
	iterator := db.db.NewIterator(&util.Range{Start: TableTransferStart[:], Limit: TableTransferLimit[:]}, nil)
	defer iterator.Release()
	receipts := []transfer.TransferReceipt{}
	for iterator.Next() {
		value := iterator.Value()
		receipt := transfer.TransferReceipt{}
		if err := json.Unmarshal(value, &receipt); err != nil {
			return receipts, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, iterator.Error()
}

func (db *dbStorage) UpdateTransferReceipt(updateReceipt transfer.UpdateReceipt) error {
	txHash, err := txHashToBytes(updateReceipt.TxHash)
	if err != nil {
		return err
	}
	receiptBytes, err := db.db.Get(append(TableTransfer[:], txHash...), nil)
	if err != nil {
		return err
	}
	receipt := transfer.TransferReceipt{}
	if err := json.Unmarshal(receiptBytes, &receipt); err != nil {
		return err
	}
	updateReceipt.Update(&receipt)
	updatedReceiptBytes, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	return db.db.Put(append(TableTransfer[:], txHash...), updatedReceiptBytes, nil)
}

func txHashToBytes(txHash string) ([]byte, error) {
	if len(txHash) > 2 && txHash[:2] == "0x" {
		txHash = txHash[2:]
	}
	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}
	if len(txHashBytes) != 32 {
		return nil, fmt.Errorf("unexpected tx hash length: %d", len(txHashBytes))
	}
	return txHashBytes, nil
}
