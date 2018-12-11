package swap

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
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

// The SwapReceipt contains the swap details and the status.
type SwapReceipt struct {
	ID            SwapID `json:"id"`
	SendToken     string `json:"sendToken"`
	ReceiveToken  string `json:"receiveToken"`
	SendAmount    string `json:"sendAmount"`
	ReceiveAmount string `json:"receiveAmount"`
	Timestamp     int64  `json:"timestamp"`
	Status        int    `json:"status"`
}

// NewSwapReceipt returns a SwapReceipt from a swapBlob.
func NewSwapReceipt(blob SwapBlob) SwapReceipt {
	return SwapReceipt{blob.ID, blob.SendToken, blob.ReceiveToken, blob.SendAmount, blob.ReceiveAmount, time.Now().Unix(), 1}
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
	ID           SwapID `json:"id"`
	SendToken    string `json:"sendToken"`
	ReceiveToken string `json:"receiveToken"`

	// SendAmount and ReceiveAmount are decimal strings.
	SendFee              string `json:"sendFee"`
	SendAmount           string `json:"sendAmount"`
	ReceiveFee           string `json:"receiveFee"`
	ReceiveAmount        string `json:"receiveAmount"`
	MinimumReceiveAmount string `json:"minimumReceiveAmount,omitempty"`

	SendTo              string `json:"sendTo"`
	ReceiveFrom         string `json:"receiveFrom"`
	TimeLock            int64  `json:"timeLock"`
	SecretHash          string `json:"secretHash"`
	ShouldInitiateFirst bool   `json:"shouldInitiateFirst"`

	Delay            bool            `json:"delayed,omitempty"`
	DelayInfo        json.RawMessage `json:"delayInfo,omitempty"`
	DelayCallbackURL string          `json:"delayCallbackUrl,omitempty"`

	BrokerFee              int64  `json:"brokerFee"` // should be between 0 and 100
	BrokerSendTokenAddr    string `json:"brokerSendTokenAddr"`
	BrokerReceiveTokenAddr string `json:"brokerReceiveTokenAddr"`

	Password string `json:"password,omitempty"`
}

type ReceiptQuery struct {
	Responder chan<- map[SwapID]SwapReceipt
}
