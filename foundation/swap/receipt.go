package swap

import (
	"encoding/json"
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

// The SwapReceipt contains the swap details and the status.
type SwapReceipt struct {
	ID            SwapID               `json:"id"`
	SendToken     blockchain.TokenName `json:"sendToken"`
	ReceiveToken  blockchain.TokenName `json:"receiveToken"`
	SendAmount    string               `json:"sendAmount"`
	ReceiveAmount string               `json:"receiveAmount"`
	SendCost      blockchain.CostBlob  `json:"sendCost"`
	ReceiveCost   blockchain.CostBlob  `json:"receiveCost"`
	Timestamp     int64                `json:"timestamp"`
	TimeLock      int64                `json:"timeLock"`
	Status        int                  `json:"status"`
	Delay         bool                 `json:"delay"`
	DelayInfo     json.RawMessage      `json:"delayInfo,omitempty"`
	Active        bool                 `json:"active"`
	PasswordHash  string               `json:"passwordHash,omitempty"`
}

// NewSwapReceipt returns a SwapReceipt from a swapBlob.
func NewSwapReceipt(blob SwapBlob) SwapReceipt {
	return SwapReceipt{
		ID:            blob.ID,
		SendToken:     blob.SendToken,
		ReceiveToken:  blob.ReceiveToken,
		SendAmount:    blob.SendAmount,
		ReceiveAmount: blob.ReceiveAmount,
		SendCost:      blockchain.CostBlob{},
		ReceiveCost:   blockchain.CostBlob{},
		Timestamp:     time.Now().Unix(),
		TimeLock:      blob.TimeLock,
		Status:        0,
		Delay:         blob.Delay,
		DelayInfo:     blob.DelayInfo,
		Active:        true,
		PasswordHash:  blob.PasswordHash,
	}
}

type ReceiptUpdate struct {
	ID     SwapID
	Update func(receipt *SwapReceipt)
}

func NewReceiptUpdate(id SwapID, update func(receipt *SwapReceipt)) ReceiptUpdate {
	return ReceiptUpdate{id, update}
}
