package swap

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

// A SwapID uniquely identifies a Swap that is being executed.
type SwapID string

func RandomID() SwapID {
	id := [32]byte{}
	rand.Read(id[:])
	return SwapID(base64.StdEncoding.EncodeToString(id[:]))
}

const ExpiryUnit = int64(2 * 60 * 60)

// The SwapReceipt contains the swap details and the status.
type SwapReceipt struct {
	ID            SwapID              `json:"id"`
	SendToken     string              `json:"sendToken"`
	ReceiveToken  string              `json:"receiveToken"`
	SendAmount    string              `json:"sendAmount"`
	ReceiveAmount string              `json:"receiveAmount"`
	SendCost      blockchain.CostBlob `json:"sendCost"`
	ReceiveCost   blockchain.CostBlob `json:"receiveCost"`
	Timestamp     int64               `json:"timestamp"`
	TimeLock      int64               `json:"timeLock"`
	Status        int                 `json:"status"`
	Delay         bool                `json:"delay"`
	DelayInfo     json.RawMessage     `json:"delayInfo,omitempty"`
}

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
	}
}

// A Swap stores all of the information required to execute an atomic swap.
type Swap struct {
	ID              SwapID
	Token           blockchain.Token
	Value           *big.Int
	Fee             *big.Int
	BrokerFee       *big.Int
	SecretHash      [32]byte
	TimeLock        int64
	SpendingAddress string
	FundingAddress  string
	BrokerAddress   string
}

// A SwapBlob is used to encode a Swap for storage and transmission.
type SwapBlob struct {
	ID           SwapID `json:"id,omitempty"`
	SendToken    string `json:"sendToken"`
	ReceiveToken string `json:"receiveToken"`

	// SendAmount and ReceiveAmount are decimal strings.
	SendFee              string `json:"sendFee,omitempty"`
	SendAmount           string `json:"sendAmount"`
	ReceiveFee           string `json:"receiveFee,omitempty"`
	ReceiveAmount        string `json:"receiveAmount"`
	MinimumReceiveAmount string `json:"minimumReceiveAmount,omitempty"`

	SendTo              string `json:"sendTo"`
	ReceiveFrom         string `json:"receiveFrom"`
	TimeLock            int64  `json:"timeLock"`
	SecretHash          string `json:"secretHash"`
	ShouldInitiateFirst bool   `json:"shouldInitiateFirst"`

	Delay            bool            `json:"delay,omitempty"`
	DelayInfo        json.RawMessage `json:"delayInfo,omitempty"`
	DelayCallbackURL string          `json:"delayCallbackUrl,omitempty"`

	BrokerFee              int64  `json:"brokerFee,omitempty"` // in BIPs or (1/10000)
	BrokerSendTokenAddr    string `json:"brokerSendTokenAddr,omitempty"`
	BrokerReceiveTokenAddr string `json:"brokerReceiveTokenAddr,omitempty"`

	ResponseURL string `json:"responseURL,omitempty"`
	Password    string `json:"password,omitempty"`
}

type ReceiptUpdate struct {
	ID     SwapID
	Update func(receipt *SwapReceipt)
}

func NewReceiptUpdate(id SwapID, update func(receipt *SwapReceipt)) ReceiptUpdate {
	return ReceiptUpdate{id, update}
}
