package swap

import "github.com/republicprotocol/renex-swapper-go/domain/order"

// Swapper is the interface for an atomic swapper object
type Swapper interface {
	NewSwap(orderID order.ID, req Request) (Swap, error)
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

func (swapper *swapper) NewSwap(orderID order.ID, req Request) (Swap, error) {
	personalAtom, foreignAtom, _, adapter, err := swapper.adapter.NewSwap(orderID, req)
	if err != nil {
		return nil, err
	}
	return &swap{
		personalAtom: personalAtom,
		foreignAtom:  foreignAtom,
		Adapter:      adapter,
	}, nil
}
