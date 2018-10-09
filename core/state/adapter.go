package state

import "github.com/republicprotocol/swapperd/service/logger"

type Adapter interface {
	logger.Logger
	Store
}

type Store interface {
	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error
	Delete([]byte) error
}
