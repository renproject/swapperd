package renex

import (
	"fmt"
	"log"

	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type renex struct {
	Adapter
	swapStatuses map[[32]byte]bool
	notifyCh     chan struct{}
	doneCh       chan struct{}
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
		Adapter:      adapter,
		swapStatuses: map[[32]byte]bool{},
		notifyCh:     make(chan struct{}, 1),
		doneCh:       make(chan struct{}, 1),
	}
}

// Run runs the watch object on the given order id
func (renex *renex) Start() <-chan error {
	errs := make(chan error)
	log.Println("Starting the watcher......")
	go func() {
		defer close(errs)
		defer log.Println("Stopping the watcher......")
		for {
			select {
			case <-renex.doneCh:
				return
			case <-renex.notifyCh:
				swaps, err := renex.ExecutableSwaps()
				if err != nil {
					errs <- err
					continue
				}
				if len(swaps) < 1000 {
					renex.SwapMultiple(swaps, errs)
					continue
				}
				renex.SwapMultiple(swaps[:1000], errs)
			}
		}
	}()
	return errs
}

func (renex *renex) SwapMultiple(swaps [][32]byte, errs chan error) {
	for i := range swaps {
		go func(i int) {
			if renex.swapStatuses[swaps[i]] {
				return
			}
			renex.swapStatuses[swaps[i]] = true
			if err := renex.Swap(swaps[i]); err != nil {
				errs <- err
			}
			if err := renex.DeleteIfRedeemedOrExpired(swaps[i]); err != nil {
				errs <- err
			}
			renex.swapStatuses[swaps[i]] = false
		}(i)
	}
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

func (renex *renex) initiate(orderID [32]byte) error {
	renex.LogInfo(orderID, "starting the atomic swap")

	if err := renex.PutStatus(orderID, swap.StatusPending); err != nil {
		return err
	}

	renex.LogInfo(orderID, "started the atomic swap")
	return nil
}

func (renex *renex) getMatch(orderID [32]byte) error {
	renex.LogInfo(orderID, "waiting for the match to be found")

	timeStamp, err := renex.SwapTimestamp(orderID)
	if err != nil {
		return err
	}

	match, err := renex.GetOrderMatch(orderID, timeStamp+48*60*60)
	if err != nil {
		renex.LogInfo(orderID, "deleting expired or unauthorized order")
		return renex.PutStatus(orderID, swap.StatusExpired)
	}

	if err := renex.PutMatch(orderID, match); err != nil {
		return err
	}

	if err := renex.PutStatus(orderID, swap.StatusMatched); err != nil {
		return err
	}

	renex.LogInfo(orderID, fmt.Sprintf("<----------> (%s)", order.Fmt(match.ForeignOrderID())))
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
