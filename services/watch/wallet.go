package watch

import "github.com/republicprotocol/atom-go/services/swap"

// Wallet is an interface for the Atom Wallet Contract
type Wallet interface {
	WaitForMatch([32]byte) ([32]byte, error)

	GetMatch([32]byte, [32]byte) (swap.OrderMatch, error)
}
