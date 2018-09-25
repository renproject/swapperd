package swap

import (
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type swapAdapter struct {
	personalAtom swap.Atom
	foreignAtom  swap.Atom
	req          swapDomain.Request
}

func (swap *swapAdapter) Initate() error {
	return swap.personalAtom.Initiate()
}

func (swap *swapAdapter) Audit() error {
	return swap.foreignAtom.Audit()
}

func (swap *swapAdapter) Redeem(secret [32]byte) error {
	return swap.foreignAtom.Redeem(secret)
}

func (swap *swapAdapter) AuditSecret() ([32]byte, error) {
	return swap.personalAtom.AuditSecret()
}

func (swap *swapAdapter) Refund() error {
	return swap.personalAtom.Refund()
}
