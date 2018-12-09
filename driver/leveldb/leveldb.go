package leveldb

import (
	"path"

	"github.com/syndtr/goleveldb/leveldb"
)

func NewStore(homeDir, network string) (*leveldb.DB, error) {
	return leveldb.OpenFile(path.Join(homeDir, "db", network), nil)
}
