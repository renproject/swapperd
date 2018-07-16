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
	if err != nil {
		return []byte{}, err
	}

	value, err := db.Get(key, nil)
	if err != nil {
		return []byte{}, err
	}

	db.Close()
	return value, nil
}

func (ldb *ldbStore) Write(key []byte, value []byte) error {
	db, err := leveldb.OpenFile(ldb.path, nil)

	if err != nil {
		return err
	}

	err := db.Put(key, value, nil)

	if err != nil {
		return err
	}

	db.Close()
	return nil
}
