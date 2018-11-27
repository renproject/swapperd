package leveldb

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func NewStore() (*leveldb.DB, error) {
	return leveldb.OpenFile(buildDBPath(), nil)
}

func buildDBPath() string {
	unix := os.Getenv("HOME")
	if unix != "" {
		return unix + "/.swapperd/db"
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return "C:\\windows\\system32\\config\\systemprofile\\swapperd\\db"
	}
	panic("unknown Operating System")
}
