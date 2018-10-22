package foundation

import (
	"math/big"
)

type SwapID [32]byte

type SwapStatus struct {
	ID     SwapID `json:"id"`
	Status Status `json:"status"`
}

type Swap struct {
	ID                 SwapID   `json:"id"`
	Secret             [32]byte `json:"secret"`
	SecretHash         [32]byte `json:"secretHash"`
	TimeLock           int64    `json:"timeLock"`
	SendToAddress      string   `json:"sendToAddress"`
	ReceiveFromAddress string   `json:"receiveFromAddress"`
	SendValue          *big.Int `json:"sendValue"`
	ReceiveValue       *big.Int `json:"receiveValue"`
	SendToken          Token    `json:"sendToken"`
	ReceiveToken       Token    `json:"receiveToken"`
	IsFirst            bool     `json:"isFirst"`
}

type SwapTry struct {
	ID              SwapID
	Token           Token
	Value           *big.Int
	SecretHash      [32]byte
	TimeLock        int64
	SpendingAddress string
	FundingAddress  string
}
