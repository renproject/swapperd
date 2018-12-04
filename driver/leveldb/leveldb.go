package leveldb

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func NewStore(network string) (*leveldb.DB, error) {
	return leveldb.OpenFile(buildDBPath(network), nil)
}

func buildDBPath(network string) string {
	unix := os.Getenv("HOME")
	if unix != "" {
		return unix + "/.swapperd/db/" + network
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return "C:\\windows\\system32\\config\\systemprofile\\swapperd\\db\\" + network
	}
	panic("unknown Operating System")
}
