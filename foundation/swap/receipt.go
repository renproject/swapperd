package swap

import (
	"encoding/json"

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
