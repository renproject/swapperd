package renex

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type renex struct {
	Adapter
	notifyCh chan struct{}
	doneCh   chan struct{}
}

type RenEx interface {
	Start() <-chan error
	Add([32]byte) error
	Status([32]byte) swap.Status
	Notify()
	Stop()
}

func NewRenEx(adapter Adapter) RenEx {
	return &renex{
		Adapter:  adapter,
		notifyCh: make(chan struct{}, 1),
		doneCh:   make(chan struct{}, 1),
	}
}

// Run runs the watch object on the given order id
func (renex *renex) Start() <-chan error {
	errs := make(chan error)
	fullsync := true
	log.Println("Starting the watcher......")
	go func() {
		defer close(errs)
		defer log.Println("Stopping the watcher......")
		for {
			select {
			case <-renex.doneCh:
				return
			case <-renex.notifyCh:
				swaps, err := renex.ExecutableSwaps(fullsync)
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
							if err := renex.Swap(swaps[i]); err != nil {
								errs <- err
								return
							}
							if renex.Status(swaps[i]) == swap.StatusRedeemed {
								renex.DeleteSwap(swaps[i])
							}
						}(i)
					}
					continue
				}
				for i := range swaps[:1000] {
					go func(i int) {
						if err := renex.Swap(swaps[i]); err != nil {
							errs <- err
							return
						}
						if renex.Status(swaps[i]) == swap.StatusRedeemed {
							renex.DeleteSwap(swaps[i])
						}
					}(i)
				}
			}
		}
	}()
	return errs
}

func (renex *renex) Add(orderID [32]byte) error {
	return renex.AddSwap(orderID)
}

func (renex *renex) Notify() {
	renex.notifyCh <- struct{}{}
}

func (renex *renex) Stop() {
	renex.doneCh <- struct{}{}
}

func (renex *renex) Swap(orderID [32]byte) error {
	if renex.Status(orderID) == swap.StatusUnknown {
		if err := renex.initiate(orderID); err != nil {
			renex.LogError(orderID, fmt.Sprintf("failed to initiate the watcher on %v", err))
			return fmt.Errorf("failed to initiate the watcher on %v", err)
		}
	} else {
		renex.LogInfo(orderID, "skipping watcher initiation")
	}

	if renex.Status(orderID) == "PENDING" {
		if err := renex.getMatch(orderID); err != nil {
			renex.LogError(orderID, fmt.Sprintf("failed to get the matching order %v", err))
			return fmt.Errorf("failed to get the matching order %v", err)
		}
	} else {
		renex.LogInfo(orderID, "skipping get match")
	}

	if renex.Status(orderID) == "MATCHED" {
		if err := renex.setInfo(orderID); err != nil {
			renex.LogError(orderID, fmt.Sprintf("failed to send address %v", err))
			return fmt.Errorf("failed to send address %v", err)
		}
	} else {
		renex.LogInfo(orderID, "skipping address submission")
	}

	if renex.Status(orderID) != "REDEEMED" && renex.Status(orderID) != "REFUNDED" && renex.Status(orderID) != swap.StatusComplained {
		if err := renex.execute(orderID); err != nil {
			renex.LogError(orderID, fmt.Sprintf("failed to execute the atomic swap %v", err))
			return fmt.Errorf("failed to execute the atomic swap %v", err)
		}
	} else {
		renex.LogInfo(orderID, "skipping address submission")
	}

	return nil
}

func (renex *renex) setInfo(orderID [32]byte) error {
	renex.LogInfo(orderID, "submitting address")
	m, err := renex.Match(orderID)
	if err != nil {
		return err
	}

	if err := renex.SendOwnerAddress(orderID, renex.GetAddress(m.ReceiveCurrency())); err != nil {
		return err
	}

	if err := renex.PutStatus(orderID, swap.StatusInfoSubmitted); err != nil {
		return err
	}

	renex.LogInfo(orderID, "submitted the address")
	return nil
}

func (renex *renex) execute(orderID [32]byte) error {
	swap, err := renex.NewSwap(orderID)
	if err != nil {
		return err
	}
	return swap.Execute()
}

func (renex *renex) initiate(orderID [32]byte) error {
	renex.LogInfo(orderID, "starting the atomic swap")

	if err := renex.PutOrderTimeStamp(orderID); err != nil {
		return err
	}

	if err := renex.PutStatus(orderID, "PENDING"); err != nil {
		return err
	}

	renex.LogInfo(orderID, "started the atomic swap")
	return nil
}

func (renex *renex) getMatch(orderID [32]byte) error {
	renex.LogInfo(orderID, "waiting for the match to be found")

	timeStamp, err := renex.OrderTimeStamp(orderID)
	if err != nil {
		return err
	}

	match, err := renex.GetOrderMatch(orderID, timeStamp+48*60*60)
	if err != nil {
		return err
	}

	if err := renex.PutMatch(orderID, match); err != nil {
		return err
	}

	if err := renex.PutStatus(orderID, "MATCHED"); err != nil {
		return err
	}

	renex.LogInfo(orderID, fmt.Sprintf("<----------> (%s)", order.Fmt(match.ForeignOrderID())))
	return nil
}
