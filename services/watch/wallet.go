package watch

import (
	"github.com/republicprotocol/atom-go/domains/match"
)

// Wallet is an interface for the Atom Wallet Contract
type Wallet interface {
	// TODO: Idiomatic Go requires this method to be called "Match" instead of
	// "GetMatch"
	GetMatch([32]byte) (match.Match, error)

	// TODO: These methods should be commented.
	SetMatch(match.Match) error
}
