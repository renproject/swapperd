package renex

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/co-go"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

type renex struct {
	Adapter
	manager  utils.SwapManager
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
		manager:  utils.NewSwapManager(),
		notifyCh: make(chan struct{}, 1),
		doneCh:   make(chan struct{}, 1),
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
					fmt.Println("Error while loading executable swaps")
					errs <- err
					continue
				}
				renex.SwapMultiple(swaps, errs)
			}
		}
	}()
	return errs
}

func (renex *renex) SwapMultiple(swaps [][32]byte, errs chan error) {
	for _, swap := range swaps {
		fmt.Println(hex.EncodeToString(swap[:]))
	}
	co.ParForAll(swaps, func(i int) {
		swap := swaps[i]
		if !renex.manager.Lock(swap) {
			return
		}
		defer renex.manager.Unlock(swap)
		if err := renex.Swap(swap); err != nil {
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
	})
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
	renex.LogInfo(orderID, "Watching RenEx Settlement for order matches")
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
	renex.LogInfo(orderID, "building swap request")
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

	sendToAddress, receiveFromAddress := renex.GetAddresses(ordMatch.ReceiveToken, ordMatch.SendToken)
	req.SendToken = ordMatch.SendToken
	req.ReceiveToken = ordMatch.ReceiveToken
	req.SendValue = ordMatch.SendValue
	req.ReceiveValue = ordMatch.ReceiveValue

	if req.SendToken > req.ReceiveToken {
		req.GoesFirst = true
		rand.Read(req.Secret[:])
		req.SecretHash = sha256.Sum256(req.Secret[:])
		req.TimeLock = time.Now().Unix() + 48*60*60
	} else {
		req.GoesFirst = false
	}

	renex.LogInfo(req.UID, "communicating swap details")
	if err := renex.SendSwapDetails(req.UID, SwapDetails{
		SecretHash:         req.SecretHash,
		TimeLock:           req.TimeLock,
		SendToAddress:      sendToAddress,
		ReceiveFromAddress: receiveFromAddress,
	}); err != nil {
		return req, err
	}

	foreignDetails, err := renex.ReceiveSwapDetails(ordMatch.ForeignOrderID, timeStamp+48*60*60)
	if err != nil {
		return req, err
	}
	renex.LogInfo(req.UID, "communication successful")

	req.SendToAddress = foreignDetails.SendToAddress
	req.ReceiveFromAddress = foreignDetails.ReceiveFromAddress
	if !req.GoesFirst {
		req.SecretHash = foreignDetails.SecretHash
		req.TimeLock = foreignDetails.TimeLock
	}

	renex.PrintSwapRequest(req)
	return req, nil
}
