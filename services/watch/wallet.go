package watch

import (
	"github.com/republicprotocol/atom-go/domains/match"
)

// Wallet is an interface for the Atom Wallet Contract
type Wallet interface {
	GetMatch([32]byte) (match.Match, error)
	SetMatch(match.Match) error
}
