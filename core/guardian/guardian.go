package guardian

import (
	"fmt"
	"time"

	"github.com/republicprotocol/co-go"
	swapDomain "github.com/republicprotocol/swapperd/domain/swap"
)

var (
	ErrNotRefundable = fmt.Errorf("Trying to refund a non refundable order")
)

func ErrAtomBuildFailed(err error) error {
	return fmt.Errorf("Failed to build the atom: %v", err)
}

func ErrRefundAfterRedeem(err error) error {
	return fmt.Errorf("Trying to an atomic swap refund after redeem: %v", err)
}

type Guardian interface {
	Run(chan<- error)
}

type guardian struct {
	Adapter
}

func NewGuardian(adapter Adapter) Guardian {
	return &guardian{
		Adapter: adapter,
	}
}

func (g *guardian) Run(errCh chan<- error) {
	for {
		swaps := g.ActiveSwaps()
		co.ParForAll(swaps, func(i int) {
			swap := swaps[i]
			details := g.SwapDetails(swap)
			if details.Status == swapDomain.StatusOpen && time.Now().Unix() > details.TimeStamp+48*60*60 {
				if err := g.PutStatus(swap, swapDomain.StatusExpired); err != nil {
					errCh <- err
					return
				}
			}
			if details.Status == swapDomain.StatusConfirmed && time.Now().Unix() > details.Request.TimeLock {
				if err := g.Refund(swap); err != nil {
					if err == ErrNotRefundable {
						return
					}
					errCh <- err
					return
				}
				if err := g.PutStatus(swap, swapDomain.StatusExpired); err != nil {
					errCh <- err
					return
				}
			}
			if err := g.DeleteIfExpired(swap); err != nil {
				errCh <- err
				return
			}
		})
		time.Sleep(1 * time.Minute)
	}
}
