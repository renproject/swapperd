package swap

import "math/big"

type Atom interface {
	Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit(hash [32]byte, to []byte, value *big.Int, delay int64) error
	GetSecretHash() [32]byte
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	PriorityCode() uint32
	From() []byte
}
