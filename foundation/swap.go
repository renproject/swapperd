package foundation

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"
)

const ExpiryUnit = int64(2 * 60 * 60)

// A SwapID uniquely identifies a Swap that is being executed.
type SwapID string

// RandomID create a random SwapID.
func RandomID() SwapID {
	id := [32]byte{}
	rand.Read(id[:])
	return SwapID(base64.StdEncoding.EncodeToString(id[:]))
}

// NewSwapStatus returns the `SwapStatus` from the given SwapBlob.
func NewSwapStatus(blob SwapBlob) SwapStatus {
	return SwapStatus{blob.ID, blob.SendToken, blob.ReceiveToken, blob.SendAmount, blob.ReceiveAmount, time.Now().Unix(), 1}
}

// A Swap stores all of the information required to execute an atomic swap.
type Swap struct {
	ID              SwapID
	Token           Token
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
	ID           SwapID `json:"id"`
	SendToken    string `json:"sendToken"`
	ReceiveToken string `json:"receiveToken"`

	// SendAmount and ReceiveAmount are decimal strings.
	SendFee              string `json:"sendFee"`
	SendAmount           string `json:"sendAmount"`
	ReceiveFee           string `json:"receiveFee"`
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

	BrokerFee              int64  `json:"brokerFee"` // should be between 0 and 100
	BrokerSendTokenAddr    string `json:"brokerSendTokenAddr"`
	BrokerReceiveTokenAddr string `json:"brokerReceiveTokenAddr"`
}

type SwapRequest struct {
	SwapBlob

	Secret   [32]byte `json:"secret"`
	Password string   `json:"password"`
}

func NewSwapRequest(swapBlob SwapBlob, secret [32]byte, password string) SwapRequest {
	return SwapRequest{swapBlob, secret, password}
}

type SwapResult struct {
	ID      SwapID
	Success bool
}

func NewSwapResult(id SwapID, success bool) SwapResult {
	return SwapResult{id, success}
}
