package swap

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"math/rand"
	"reflect"

	"github.com/renproject/swapperd/foundation/blockchain"
)

const ExpiryUnit = int64(2 * 60 * 60)

// TODO: Rename to ID
// A SwapID uniquely identifies a Swap that is being executed.
type SwapID string

// Generate is used to create random values for testing
func (SwapID) Generate(rand *rand.Rand, size int) reflect.Value {
	id := [32]byte{}
	rand.Read(id[:])
	return reflect.ValueOf(SwapID(base64.StdEncoding.EncodeToString(id[:])))
}

// A Swap stores all of the information required to execute an atomic swap.
type Swap struct {
	ID              SwapID
	Token           blockchain.Token
	Value           *big.Int
	BrokerFee       *big.Int
	SecretHash      [32]byte
	TimeLock        int64
	SpendingAddress string
	WithdrawAddress string
	FundingAddress  string
	BrokerAddress   string
	Speed           blockchain.TxExecutionSpeed
}

// A SwapBlob is used to encode a Swap for storage and transmission.
type SwapBlob struct {
	ID           SwapID               `json:"id,omitempty"`
	SendToken    blockchain.TokenName `json:"sendToken"`
	ReceiveToken blockchain.TokenName `json:"receiveToken"`

	// SendAmount and ReceiveAmount are decimal strings.
	SendAmount           string                      `json:"sendAmount"`
	ReceiveAmount        string                      `json:"receiveAmount"`
	Speed                blockchain.TxExecutionSpeed `json:"speed"`
	MinimumReceiveAmount string                      `json:"minimumReceiveAmount,omitempty"`

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

	WithdrawAddress string `json:"withdrawAddress,omitempty"`
	ResponseURL     string `json:"responseURL,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordHash    string `json:"passwordHash,omitempty"`
}
