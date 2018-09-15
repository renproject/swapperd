package watch

import (
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/store"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

// Wallet is an interface for the Atom Wallet Contract
type Adapter interface {
	// TODO: Idiomatic Go requires this method to be called "Match" instead of
	// "GetMatch"
	swap.SwapAdapter
	BuildAtoms(store.State, match.Match) (swap.Atom, swap.Atom, error)
	CheckForMatch(order.ID, bool) (match.Match, error)
}
