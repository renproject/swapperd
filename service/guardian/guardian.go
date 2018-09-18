package guardian

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

var ErrSwapRedeemed = fmt.Errorf("Swap Redeemed")

type Guardian interface {
	Start() <-chan error
	Notify()
	Stop()
}

type guardian struct {
	Adapter
	notifyCh chan struct{}
	doneCh   chan struct{}
}

func NewGuardian(adapter Adapter) Guardian {
	return &guardian{
		Adapter:  adapter,
		notifyCh: make(chan struct{}, 1),
		doneCh:   make(chan struct{}, 1),
	}
}

func (g *guardian) Start() <-chan error {
	errs := make(chan error)
	log.Println("Starting the guardian......")
	go func() {
		defer log.Println("Ending the guardian......")
		for {
			select {
			case <-g.doneCh:
				return
			case <-g.notifyCh:
				swaps, err := g.RefundableSwaps()
				if err != nil {
					if err == ErrSwapRedeemed {
						continue
					}
					errs <- err
					return
				}
				if len(swaps) < 1000 {
					for i := range swaps {
						go func(i int) {
							if err := g.Refund(swaps[i]); err != nil {
								errs <- err
								return
							}
							if g.Status(swaps[i]) == swap.StatusRefunded {
								g.DeleteSwap(swaps[i])
							}
						}(i)
					}
					continue
				}
				for i := range swaps[:1000] {
					go func(i int) {
						if err := g.Refund(swaps[i]); err != nil {
							errs <- err
							return
						}
						if g.Status(swaps[i]) == swap.StatusRefunded {
							g.DeleteSwap(swaps[i])
						}
					}(i)
				}
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
