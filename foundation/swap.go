package foundation

import (
	"math/big"
)

type SwapID [32]byte

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
