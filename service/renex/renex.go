package renex

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

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

// TODO: Change Start to Run
// Run runs the watch object on the given order id
func (renex *renex) Start() <-chan error {
	errs := make(chan error)
	log.Println("Running the watcher......")
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
				// TODO: Document the limitation
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
		// TODO: Use co library
		go func(i int) {
			// TODO: Fix data race issues
			if renex.swapStatuses[swaps[i]] {
				return
			}
			renex.swapStatuses[swaps[i]] = true
			defer func() { renex.swapStatuses[swaps[i]] = false }()
			if err := renex.Swap(swaps[i]); err != nil {
				select {
				case _, ok := <-renex.doneCh:
					if !ok {
						return
					}
				case errs <- err:
				}
			}
			if err := renex.DeleteIfRedeemedOrExpired(swaps[i]); err != nil {
				select {
				case _, ok := <-renex.doneCh:
					if !ok {
						return
					}
				case errs <- err:
				}
			}
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
	if renex.Status(orderID) == swap.StatusOpen {
		req, err := renex.buildRequest(orderID)
		if err != nil {
			return err
		}
		if err := renex.PutSwapRequest(orderID, req); err != nil {
			return err
		}
		if err := renex.PutStatus(orderID, swap.StatusConfirmed); err != nil {
			return err
		}
	}

	if renex.Status(orderID) == swap.StatusConfirmed {
		req, err := renex.SwapRequest(orderID)
		if err != nil {
			return err
		}

		swapInst, err := renex.NewSwap(req)
		if err != nil {
			return err
		}
		if err := swapInst.Execute(); err != nil {
			return err
		}
		if err := renex.PutStatus(orderID, swap.StatusSettled); err != nil {
			return err
		}
	}
	return nil
}

func (renex *renex) buildRequest(orderID [32]byte) (swap.Request, error) {
	fmt.Println("Building Request.... ")
	req := swap.Request{}

	// TODO: Change it to AddedAtTimestamp
	timeStamp, err := renex.AddTimestamp(orderID)
	if err != nil {
		return req, err
	}

	req.UID = orderID
	ordMatch, err := renex.GetOrderMatch(orderID, timeStamp+48*60*60)
	if err != nil {
		return req, err
	}
	fmt.Println(ordMatch)

	sendToAddress, receiveFromAddress := renex.GetAddresses(ordMatch.ReceiveToken, ordMatch.SendToken)
	req.SendToken = ordMatch.SendToken
	req.ReceiveToken = ordMatch.ReceiveToken
	req.SendValue = ordMatch.SendValue
	req.ReceiveValue = ordMatch.ReceiveValue

	fmt.Println()

	if req.SendToken > req.ReceiveToken {
		req.GoesFirst = true
		rand.Read(req.Secret[:])
		req.SecretHash = sha256.Sum256(req.Secret[:])
		req.TimeLock = time.Now().Unix() + 48*60*60
	} else {
		req.GoesFirst = false
	}

	fmt.Println("Communicating swap details")
	if err := renex.SendSwapDetails(req.UID, SwapDetails{
		SecretHash:         req.SecretHash,
		TimeLock:           req.TimeLock,
		SendToAddress:      sendToAddress,
		ReceiveFromAddress: receiveFromAddress,
	}); err != nil {
		return req, err
	}
	fmt.Println("Swap details sent")

	foreignDetails, err := renex.ReceiveSwapDetails(ordMatch.ForeignOrderID, timeStamp+48*60*60)
	if err != nil {
		return req, err
	}
	fmt.Println(foreignDetails)

	req.SendToAddress = foreignDetails.SendToAddress
	req.ReceiveFromAddress = foreignDetails.ReceiveFromAddress
	if !req.GoesFirst {
		req.SecretHash = foreignDetails.SecretHash
		req.TimeLock = foreignDetails.TimeLock
	}

	renex.PrintSwapRequest(req)
	return req, nil
}
