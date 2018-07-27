package swap

import (
	"math/big"
)

type Atom interface {
	Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit() ([32]byte, []byte, *big.Int, int64, error)
	WaitForCounterRedemption() error
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	GetFromAddress() ([]byte, error)
	PriorityCode() uint32
}
