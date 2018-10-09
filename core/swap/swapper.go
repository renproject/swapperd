package swap

import "github.com/republicprotocol/swapperd/domain/swap"

// Swapper is the interface for an atomic swapper object
type Swapper interface {
	NewSwap(req swap.Request) (Swap, error)
}

type swapper struct {
	adapter SwapperAdapter
}

// NewSwapper returns a new Swapper instance
func NewSwapper(adapter SwapperAdapter) Swapper {
	return &swapper{
		adapter: adapter,
	}
}

func (swapper *swapper) NewSwap(req swap.Request) (Swap, error) {
	personalAtom, foreignAtom, adapter, err := swapper.adapter.NewSwap(req)
	if err != nil {
		return nil, err
	}
	return &swapExec{
		req:          req,
		personalAtom: personalAtom,
		foreignAtom:  foreignAtom,
		Adapter:      adapter,
	}, nil
}
