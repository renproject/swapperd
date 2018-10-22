package leveldb

import (
	"os"

	"github.com/republicprotocol/swapperd/adapter/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewStore(homeDir string) (storage.Store, error) {
	db, err := leveldb.OpenFile(buildDBPath(homeDir), nil)
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

func buildDBPath(homeDir string) string {
	unix := os.Getenv("HOME")
	if unix != "" {
		return homeDir + "/db"
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return homeDir + "\\db"
	}
	panic("unknown Operating System")
}
