package watch

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/store"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type watch struct {
	adapter  Adapter
	state    store.State
	notifyCh chan struct{}
	doneCh   chan struct{}
}

type Watch interface {
	Start() <-chan error
	Add([32]byte) error
	Status([32]byte) string
	Notify()
	Stop()
}

func NewWatch(adapter Adapter, state store.State) Watch {
	return &watch{
		adapter:  adapter,
		state:    state,
		notifyCh: make(chan struct{}, 1),
		doneCh:   make(chan struct{}, 1),
	}
}

// Run runs the watch object on the given order id
func (watch *watch) Start() <-chan error {
	errs := make(chan error)
	fullsync := true
	log.Println("Starting the watcher......")
	go func() {
		defer close(errs)
		defer log.Println("Stopping the watcher......")
		for {
			select {
			case <-watch.doneCh:
				return
			case <-watch.notifyCh:
				swaps, err := watch.state.ExecutableSwaps(fullsync)
				if fullsync {
					fullsync = false
				}
				if err != nil {
					errs <- err
					continue
				}
				if len(swaps) < 1000 {
					for i := range swaps {
						go func(i int) {
							if err := watch.Swap(swaps[i]); err != nil {
								errs <- err
								return
							}
							if watch.state.Status(swaps[i]) == swap.StatusRedeemed {
								watch.state.DeleteSwap(swaps[i])
							}
						}(i)
					}
					continue
				}
				for i := range swaps[:1000] {
					go func(i int) {
						if err := watch.Swap(swaps[i]); err != nil {
							errs <- err
							return
						}
						if watch.state.Status(swaps[i]) == swap.StatusRedeemed {
							watch.state.DeleteSwap(swaps[i])
						}
					}(i)
				}
			}
		}
	}()
	return errs
}

func (watch *watch) Add(orderID [32]byte) error {
	return watch.state.AddSwap(orderID)
}

func (watch *watch) Status(orderID [32]byte) string {
	return watch.state.Status(orderID)
}

func (watch *watch) Notify() {
	watch.notifyCh <- struct{}{}
}

func (watch *watch) Stop() {
	watch.doneCh <- struct{}{}
}

func (watch *watch) Swap(orderID [32]byte) error {
	if watch.state.Status(orderID) == "UNKNOWN" {
		if err := watch.initiate(orderID); err != nil {
			watch.adapter.LogError(orderID, fmt.Sprintf("failed to initiate the watcher on %v", err))
			return fmt.Errorf("failed to initiate the watcher on %v", err)
		}
	} else {
		watch.adapter.LogInfo(orderID, "skipping watcher initiation")
	}

	if watch.state.Status(orderID) == "PENDING" {
		if err := watch.getMatch(orderID); err != nil {
			watch.adapter.LogError(orderID, fmt.Sprintf("failed to get the matching order %v", err))
			return fmt.Errorf("failed to get the matching order %v", err)
		}
	} else {
		watch.adapter.LogInfo(orderID, "skipping get match")
	}

	if watch.state.Status(orderID) == "MATCHED" {
		if err := watch.setInfo(orderID); err != nil {
			watch.adapter.LogError(orderID, fmt.Sprintf("failed to send address %v", err))
			return fmt.Errorf("failed to send address %v", err)
		}
	} else {
		watch.adapter.LogInfo(orderID, "skipping address submission")
	}

	if watch.state.Status(orderID) != "REDEEMED" && watch.state.Status(orderID) != "REFUNDED" && watch.state.Status(orderID) != swap.StatusComplained {
		if err := watch.execute(orderID); err != nil {
			watch.adapter.LogError(orderID, fmt.Sprintf("failed to execute the atomic swap %v", err))
			return fmt.Errorf("failed to execute the atomic swap %v", err)
		}
	} else {
		watch.adapter.LogInfo(orderID, "skipping address submission")
	}

	return nil
}

func (watch *watch) setInfo(orderID [32]byte) error {
	watch.adapter.LogInfo(orderID, "submitting address")
	m, err := watch.state.Match(orderID)
	if err != nil {
		return err
	}

	_, foreignAtom, err := watch.adapter.BuildAtoms(watch.state, m)
	if err != nil {
		return err
	}

	addr, err := foreignAtom.GetFromAddress()
	if err != nil {
		return err
	}

	if err := watch.adapter.SendOwnerAddress(orderID, addr); err != nil {
		return err
	}

	if err := watch.state.PutStatus(orderID, swap.StatusInfoSubmitted); err != nil {
		return err
	}

	watch.adapter.LogInfo(orderID, "submitted the address")
	return nil
}

func (watch *watch) execute(orderID [32]byte) error {
	m, err := watch.state.Match(orderID)
	if err != nil {
		return err
	}

	personalAtom, foreignAtom, err := watch.adapter.BuildAtoms(watch.state, m)
	if err != nil {
		return err
	}

	atomicSwap := swap.NewSwap(personalAtom, foreignAtom, m, watch.adapter, watch.state)
	return atomicSwap.Execute()
}

func (watch *watch) initiate(orderID [32]byte) error {
	watch.adapter.LogInfo(orderID, "starting the atomic swap")
	err := watch.state.PutStatus(orderID, "PENDING")
	if err != nil {
		return err
	}
	watch.adapter.LogInfo(orderID, "started the atomic swap")
	return nil
}

func (watch *watch) getMatch(orderID [32]byte) error {
	watch.adapter.LogInfo(orderID, "waiting for the match to be found")
	match, err := watch.adapter.CheckForMatch(orderID, true)
	if err != nil {
		return err
	}

	err = watch.state.PutMatch(orderID, match)
	if err != nil {
		return err
	}

	err = watch.state.PutStatus(orderID, "MATCHED")
	if err != nil {
		return err
	}

	watch.adapter.LogInfo(orderID, fmt.Sprintf("<----------> (%s)", order.Fmt(match.ForeignOrderID())))
	return nil
}
