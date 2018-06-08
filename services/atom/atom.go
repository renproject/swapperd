package atom

import (
	"math/big"
)

/*
 * Steps in our Atomic Swap:
 *
 * 0. (A) and (B) share addresses, (A) creates HASH from SECRET
 * 1. (A) calls Initiate(HASH, details) to ADDR1, gives (HASH, ADDR1) to (B)
 * 2. (B) calls ADDR1.Audit(HASH, details)
 * 3. (B) calls Initiate(HASH) to ADDR2, gives (ADDR2) to (A)
 * 4. (A) calls ADDR2.Audit(HASH, details)
 * 5. (A) calls ADDR2.Redeem(SECRET)
 * 6. (B) calls ADDR2.AuditSecret(), retrieving SECRET
 * 7. (B) calls ADDR1.Redeem(SECRET)
 */

// Atom is the interface defining the Atomic Swap Interface
type Atom interface {
	Initiate(hash [32]byte, value *big.Int, expiry int64) error
	Redeem(secret [32]byte) error
	Refund() error
	Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error)
	AuditSecret() (secret [32]byte, err error)
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	From() []byte
	PriorityCode() int64
}
