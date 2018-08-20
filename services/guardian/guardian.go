package guardian

import (
	"fmt"
	"log"
	"time"

	co "github.com/republicprotocol/co-go"
	"github.com/republicprotocol/renex-swapper-go/adapters/atoms"
	"github.com/republicprotocol/renex-swapper-go/services/errors"
	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/republicprotocol/renex-swapper-go/services/swap"
)

var ErrSwapRedeemed = fmt.Errorf("Swap Redeemed")

type Guardian interface {
	Start() <-chan error
	Notify()
	Stop()
}

type guardian struct {
	builder  atoms.AtomBuilder
	state    store.State
	notifyCh chan struct{}
	doneCh   chan struct{}
}

func NewGuardian(builder atoms.AtomBuilder, state store.State) Guardian {
	return &guardian{
		builder:  builder,
		state:    state,
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
				swaps, err := g.state.RefundableSwaps()
				if err != nil {
					if err == ErrSwapRedeemed {
						continue
					}
					errs <- err
					return
				}
				if len(swaps) < 1000 {
					co.ParForAll(swaps, func(i int) {
						if err := g.refund(swaps[i]); err != nil {
							if err == errors.ErrNotInitiated {
								return
							}
							errs <- err
							return
						}
						if g.state.Status(swaps[i]) == swap.StatusRefunded {
							g.state.DeleteSwap(swaps[i])
						}
					})
					continue
				}
				co.ParForAll(swaps[:1000], func(i int) {
					if err := g.refund(swaps[i]); err != nil {
						if err == errors.ErrNotInitiated {
							return
						}
						errs <- err
						return
					}
					if g.state.Status(swaps[i]) == swap.StatusRefunded {
						g.state.DeleteSwap(swaps[i])
					}
				})
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

func (g *guardian) refund(orderID [32]byte) error {
	if !g.state.Complained(orderID) && !g.state.IsRedeemable(orderID) {
		return errors.ErrNotInitiated
	}

	atom, err := g.buildAtom(orderID)
	if err != nil {
		return errors.ErrAtomBuildFailed(err)
	}

	if err = g.waitForExpiry(orderID); err != nil {
		return err
	}

	if err := atom.Refund(); err != nil {
		return errors.ErrRefundAfterRedeem(err)
	}
	return nil
}

func (g *guardian) buildAtom(orderID [32]byte) (swap.Atom, error) {
	m, err := g.state.Match(orderID)
	if err != nil {
		return nil, err
	}
	atom, _, err := g.builder.BuildAtoms(g.state, m)
	return atom, err
}

func (g *guardian) waitForExpiry(orderID [32]byte) error {
	expiry, _, err := g.state.InitiateDetails(orderID)
	if err != nil {
		return err
	}

	for {
		if time.Now().Unix() >= expiry {
			if g.state.Status(orderID) == swap.StatusRedeemed {
				return ErrSwapRedeemed
			}
			return nil
		}
		time.Sleep(time.Duration(time.Now().Unix()-expiry) * time.Minute)
	}
}
