package leveldb

import (
	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/syndtr/goleveldb/leveldb"
)

type ldbStore struct {
	db *leveldb.DB
}

func NewLDBStore(path string) (store.Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &ldbStore{
		db: db,
	}, nil
}

func (ldb *ldbStore) Read(key []byte) ([]byte, error) {
	return ldb.db.Get(key, nil)
}

func (ldb *ldbStore) Write(key []byte, value []byte) error {
	return ldb.db.Put(key, value, nil)
}

func (ldb *ldbStore) Delete(key []byte) error {
	return ldb.db.Delete(key, nil)
}

func (ldb *ldbStore) Close() error {
	return ldb.db.Close()
}
