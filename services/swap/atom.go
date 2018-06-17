package swap

import "math/big"

type AtomRequester interface {
	Initiate(hash [32]byte, value *big.Int, expiry int64) error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	PriorityCode() int64
	From() []byte
}

type AtomResponder interface {
	Redeem(secret [32]byte) error
	GetSecretHash() [32]byte
	Audit(hash [32]byte, to []byte, value *big.Int, delay int64) error
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	PriorityCode() int64
}
