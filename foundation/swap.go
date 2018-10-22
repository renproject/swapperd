package foundation

import (
	"math/big"
)

type SwapID string

type SwapStatus struct {
	ID     SwapID `json:"id"`
	Status Status `json:"status"`
}

type SwapRequest struct {
	ID                  SwapID `json:"id"`
	SendToken           string `json:"sendToken"`
	ReceiveToken        string `json:"receiveToken"`
	SendAmount          string `json:"sendAmount"`    // hex
	ReceiveAmount       string `json:"receiveAmount"` //hex
	SendTo              string `json:"sendTo"`
	ReceiveFrom         string `json:"receiveFrom"`
	TimeLock            int64  `json:"timeLock"`
	SecretHash          string `json:"secretHash"`
	ShouldInitiateFirst bool   `json:"shouldInitiateFirst"`
}

type Swap struct {
	ID              SwapID
	Token           Token
	Value           *big.Int
	SecretHash      [32]byte
	TimeLock        int64
	SpendingAddress string
	FundingAddress  string
}
