package atom

import "math/big"

type ResponseAtom interface {
	Redeem(secret [32]byte) error
	GetSecretHash() [32]byte
	Audit(hash [32]byte, to []byte, value *big.Int, delay int64) error
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	PriorityCode() int64
}
