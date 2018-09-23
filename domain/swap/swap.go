package swap

import (
	"math/big"

	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

type Status string

const (
	// StatusOpen is returned when the swap is waiting for the match to be
	// found.
	StatusOpen = Status("OPEN")

	// StatusSettling is returned when the swapper is waiting for the swap to
	// finish.
	StatusSettling = Status("SETTLING")

	// StatusSettled is returned when the swap is redeemed.
	StatusSettled = Status("SETTLED")

	// StatusExpired is returned if the given swap is expred or refunded.
	StatusExpired = Status("EXPIRED")
)

// const (
// 	// StatusUnknown is returned when the swap information is not found in the
// 	// local persistent storage.
// 	StatusUnknown = Status("UNKNOWN")

// 	// StatusPending is returned when the swapper is waiting for the match to be
// 	// found.
// 	StatusPending = Status("PENDING")

// 	// StatusMatched is returned when the match is found for the givn order ID.
// 	StatusMatched = Status("MATCHED")

// 	// StatusInfoSubmitted is returned when the addresses of the traders is
// 	// submitted.
// 	StatusInfoSubmitted = Status("INFO_SUBMITTED")

// 	// StatusInitiateDetailsAcquired is returned when the swapper acquires the
// 	// initiate details either from generating them or recieving them from the
// 	// initiating trader.
// 	StatusInitiateDetailsAcquired = Status("INITIATE_DETAILS_ACQUIRED")

// 	// StatusInitiated is returned when the atomic swap is initiated.
// 	StatusInitiated = Status("INITIATED")

// 	// StatusWaitingForCounterInitiation is returned when the trader is waiting
// 	// for the counter party to initiate.
// 	StatusWaitingForCounterInitiation = Status("WAITING_FOR_COUNTER_INITIATION")

// 	// StatusRedeemDetailsAcquired is returned when the swapper acquires the
// 	// secret, either from it's local storage or from the blockchain.
// 	StatusRedeemDetailsAcquired = Status("REDEEM_DETAILS_ACQUIRED")

// 	// StatusRedeemed is returned when the atomic swap is redeemed.
// 	StatusRedeemed = Status("REDEEMED")

// 	// StatusWaitingForCounterRedemption is returned when the swapper is waiting
// 	// for the counter party initiation.
// 	StatusWaitingForCounterRedemption = Status("WAITING_FOR_COUNTER_REDEMPTION")

// 	// StatusRefunded is returned when the given atomic swap is refunded.
// 	StatusRefunded = Status("REFUNDED")

// 	// StatusComplained is returned when the swapper complains about a failed
// 	// atomic swap.
// 	StatusComplained = Status("COMPLAINED")

// 	// StatusReceivedSwapDetails is returned when the swapper receives swap
// 	// details of the counter-party.
// 	StatusReceivedSwapDetails = Status("RECEIVED_SWAP_DETAILS")

// 	// StatusSentSwapDetails is returned after the swapper sent it's swap
// 	// details.
// 	StatusSentSwapDetails = Status("SENT_SWAP_DETAILS")

// 	// StatusAudited is returned after the swap is audited.
// 	StatusAudited = Status("AUDITED")

// 	// StatusExpired is returned if the given swap is expred.
// 	StatusExpired = Status("EXPIRED")
// )

type Request struct {
	UID                [32]byte    `json:"uid"`
	Secret             [32]byte    `json:"secret"`
	SecretHash         [32]byte    `json:"secretHash"`
	TimeLock           int64       `json:"timeLock"`
	SendToAddress      string      `json:"sendToAddress"`
	ReceiveFromAddress string      `json:"receiveFromAddress"`
	SendValue          *big.Int    `json:"sendValue"`
	ReceiveValue       *big.Int    `json:"receiveValue"`
	SendToken          token.Token `json:"sendToken"`
	ReceiveToken       token.Token `json:"receiveToken"`
	GoesFirst          bool        `json:"goesFirst"`
}

// Match is the order match interface
type Match struct {
	PersonalOrderID [32]byte
	ForeignOrderID  [32]byte
	SendValue       *big.Int
	ReceiveValue    *big.Int
	SendToken       token.Token
	ReceiveToken    token.Token
}
