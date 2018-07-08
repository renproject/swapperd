package leveldb

import (
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/syndtr/goleveldb/leveldb"
)

type ldbStore struct {
	path string
}

func NewLDBStore(path string) store.Store {
	return &ldbStore{
		path: path,
	}
}

func (ldb *ldbStore) Read(key []byte) ([]byte, error) {
	db, err := leveldb.OpenFile(ldb.path, nil)
	defer db.Close()

	if err != nil {
		return []byte{}, err
	}
	return db.Get(key, nil)
}

func (ldb *ldbStore) Write(key []byte, value []byte) error {
	db, err := leveldb.OpenFile(ldb.path, nil)
	defer db.Close()

	if err != nil {
		return err
	}

	return db.Put(key, value, nil)
}
