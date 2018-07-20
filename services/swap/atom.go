package swap

import (
	"math/big"

	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/store"
)

type Atom interface {
	Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit() ([32]byte, []byte, *big.Int, int64, error)
	WaitForCounterRedemption() error
	Store(store.SwapState) error
	Restore(store.SwapState) error

	PriorityCode() uint32
	GetSecretHash() [32]byte
	GetKey() Key
}

type AtomBuilder interface {
	BuildAtoms(store.SwapState, match.Match) (Atom, Atom, error)
}
