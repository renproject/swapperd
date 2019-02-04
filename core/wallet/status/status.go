package status

import (
	"fmt"

	"github.com/renproject/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
)

type Storage interface {
	PutReceipt(receipt swap.SwapReceipt) error
	UpdateReceipt(receiptUpdate swap.ReceiptUpdate) error
}

type statuses struct {
	storage Storage
}

func New(cap int, storage Storage) tau.Task {
	return tau.New(tau.NewIO(cap), &statuses{storage})
}

func (statuses *statuses) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Receipt:
		return statuses.handleReceipt(msg)
	case ReceiptUpdate:
		return statuses.handleReceiptUpdate(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in statuses: %T", msg))
	}
}

func (statuses *statuses) handleReceipt(receipt Receipt) tau.Message {
	if err := statuses.storage.PutReceipt(swap.SwapReceipt(receipt)); err != nil {
		return tau.NewError(err)
	}
	return nil
}

func (statuses *statuses) handleReceiptUpdate(update ReceiptUpdate) tau.Message {
	if err := statuses.storage.UpdateReceipt(swap.ReceiptUpdate(update)); err != nil {
		return tau.NewError(err)
	}
	return nil
}

type Receipt swap.SwapReceipt

func (Receipt) IsMessage() {
}

func NewReceipt(blob swap.SwapBlob) Receipt {
	return Receipt(swap.NewSwapReceipt(blob))
}

type ReceiptUpdate swap.ReceiptUpdate

func (ReceiptUpdate) IsMessage() {
}

type ReceiptQuery struct {
	Responder chan<- map[swap.SwapID]swap.SwapReceipt
}

func (ReceiptQuery) IsMessage() {
}

type Bootload struct {
}

func (Bootload) IsMessage() {
}
