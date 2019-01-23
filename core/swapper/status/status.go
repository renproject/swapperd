package status

import (
	"fmt"

	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
)

type statuses struct {
	statuses map[swap.SwapID]swap.SwapReceipt
}

func New(cap int) tau.Task {
	return tau.New(tau.NewIO(cap), &statuses{map[swap.SwapID]swap.SwapReceipt{}})
}

func (statuses *statuses) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Receipt:
		return statuses.handleReceipt(msg)
	case ReceiptUpdate:
		return statuses.handleReceiptUpdate(msg)
	case ReceiptQuery:
		return statuses.handleReceiptQuery(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in transfers: %T", msg))
	}
}

func (statuses *statuses) handleReceiptQuery(msg ReceiptQuery) tau.Message {
	// TODO: update to use shallow copy
	msg.Responder <- statuses.statuses
	return nil
}

func (statuses *statuses) handleReceipt(receipt Receipt) tau.Message {
	statuses.statuses[receipt.ID] = swap.SwapReceipt(receipt)
	return nil
}

func (statuses *statuses) handleReceiptUpdate(update ReceiptUpdate) tau.Message {
	receipt := statuses.statuses[update.ID]
	update.Update(&receipt)
	statuses.statuses[update.ID] = receipt
	return nil
}

type Receipt swap.SwapReceipt

func (msg Receipt) IsMessage() {
}

type ReceiptUpdate swap.ReceiptUpdate

func (msg ReceiptUpdate) IsMessage() {
}

type ReceiptQuery struct {
	Responder chan<- map[swap.SwapID]swap.SwapReceipt
}

func (msg ReceiptQuery) IsMessage() {
}
