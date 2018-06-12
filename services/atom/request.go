package atom

import (
	"math/big"
)

type RequestAtom interface {
	Initiate(hash [32]byte, value *big.Int, expiry int64) error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	PriorityCode() int64
	From() []byte
}
