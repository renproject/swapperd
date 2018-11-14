package foundation

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
)

// A SwapID uniquely identifies a Swap that is being executed.
type SwapID string

const (
	Inactive = iota
	Initiated
	Audited
	AuditFailed
	Redeemed
	Refunded
)

const ExpiryUnit = int64(2 * 60 * 60)

// The SwapStatus indicates which phase of execution a Swap is in.
type SwapStatus struct {
	ID     SwapID `json:"id"`
	Status int    `json:"status"`
}

func RandomID() SwapID {
	id := [32]byte{}
	rand.Read(id[:])
	return SwapID(base64.StdEncoding.EncodeToString(id[:]))
}

func NewSwapStatus(id SwapID, status int) SwapStatus {
	return SwapStatus{id, status}
}

// A Swap stores all of the information required to execute an atomic swap.
type Swap struct {
	ID              SwapID
	Token           Token
	Value           *big.Int
	SecretHash      [32]byte
	TimeLock        int64
	SpendingAddress string
	FundingAddress  string
}

// A SwapBlob is used to encode a Swap for storage and transmission.
type SwapBlob struct {
	ID           SwapID `json:"id"`
	SendToken    string `json:"sendToken"`
	ReceiveToken string `json:"receiveToken"`

	// SendAmount and ReceiveAmount are decimal strings.
	SendAmount           string `json:"sendAmount"`
	ReceiveAmount        string `json:"receiveAmount"`
	MinimumReceiveAmount string `json:"minimumReceiveAmount"`

	SendTo              string `json:"sendTo"`
	ReceiveFrom         string `json:"receiveFrom"`
	TimeLock            int64  `json:"timeLock"`
	SecretHash          string `json:"secretHash"`
	ShouldInitiateFirst bool   `json:"shouldInitiateFirst"`

	Delay            bool            `json:"delayed"`
	DelayInfo        json.RawMessage `json:"delayInfo"`
	DelayCallbackURL string          `json:"delayCallbackUrl"`
}
