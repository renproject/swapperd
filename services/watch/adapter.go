package watch

import (
	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
)

// Wallet is an interface for the Atom Wallet Contract
type Adapter interface {
	// TODO: Idiomatic Go requires this method to be called "Match" instead of
	// "GetMatch"
	swap.SwapAdapter
	BuildAtoms(store.State, match.Match) (swap.Atom, swap.Atom, error)
	Match([32]byte) (match.Match, error)
}
