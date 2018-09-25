package guardian

import (
	"fmt"
	"log"
)

var (
	ErrSwapRedeemed  = fmt.Errorf("Swap Redeemed")
	ErrNonRefundable = fmt.Errorf("Trying to refund a non refundable order")
	ErrNotInitiated  = fmt.Errorf("Trying to refund a swap which is not initiated")
)

func ErrAtomBuildFailed(err error) error {
	return fmt.Errorf("Failed to build the atom: %v", err)
}

func ErrRefundAfterRedeem(err error) error {
	return fmt.Errorf("Trying to an atomic swap refund after redeem: %v", err)
}

type Guardian interface {
	Start() <-chan error
	Notify()
	Stop()
}

type guardian struct {
	Adapter
	refundStatus map[[32]byte]bool
	notifyCh     chan struct{}
	doneCh       chan struct{}
}

func NewGuardian(adapter Adapter) Guardian {
	return &guardian{
		Adapter:      adapter,
		refundStatus: map[[32]byte]bool{},
		notifyCh:     make(chan struct{}, 1),
		doneCh:       make(chan struct{}, 1),
	}
}

func (g *guardian) Start() <-chan error {
	errs := make(chan error)
	log.Println("Starting the guardian......")
	go func() {
		defer log.Println("Ending the guardian......")
		defer close(errs)
		for {
			select {
			case <-g.doneCh:
				return
			case <-g.notifyCh:
				swaps, err := g.ExpiredSwaps()
				if err != nil {
					errs <- err
					return
				}
				if len(swaps) < 1000 {
					g.RefundMultiple(swaps, errs)
					continue
				}
				g.RefundMultiple(swaps[:1000], errs)
			}
		}
	}()
	return errs
}

func (g *guardian) Notify() {
	g.notifyCh <- struct{}{}
}

func (g *guardian) Stop() {
	g.doneCh <- struct{}{}
}

func (g *guardian) RefundMultiple(swaps [][32]byte, errs chan error) {
	for i := range swaps {
		go func(i int) {
			if g.refundStatus[swaps[i]] {
				return
			}
			g.refundStatus[swaps[i]] = true
			if err := g.Refund(swaps[i]); err != nil {
				if err == ErrNotInitiated {
					return
				}
				select {
				case _, ok := <-g.doneCh:
					if !ok {
						return
					}
				case errs <- err:
				}
			}
			if err := g.DeleteIfRefunded(swaps[i]); err != nil {
				select {
				case _, ok := <-g.doneCh:
					if !ok {
						return
					}
				case errs <- err:
				}
			}
			g.refundStatus[swaps[i]] = false
			return
		}(i)
	}
}
