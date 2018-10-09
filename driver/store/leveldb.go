package store

import (
	"github.com/republicprotocol/swapperd/service/state"
	"github.com/republicprotocol/swapperd/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB(path string) (state.Store, error) {
	path = utils.BuildDBPath(path)
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{
		db: db,
	}, nil
}

func (ldb *LevelDB) Read(key []byte) ([]byte, error) {
	return ldb.db.Get(key, nil)
}

func (ldb *LevelDB) Write(key []byte, value []byte) error {
	return ldb.db.Put(key, value, nil)
}

func (ldb *LevelDB) Delete(key []byte) error {
	return ldb.db.Delete(key, nil)
}

func (ldb *LevelDB) Close() error {
	return ldb.db.Close()
}
