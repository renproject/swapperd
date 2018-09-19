package swap

import (
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

// ID is the swap ID
type ID [32]byte

type Status string

const (
	// StatusUnknown is returned when the swap information is not found in the
	// local persistent storage.
	StatusUnknown = Status("UNKNOWN")

	// StatusPending is returned when the swapper is waiting for the match to be
	// found.
	StatusPending = Status("PENDING")

	// StatusMatched is returned when the match is found for the givn order ID.
	StatusMatched = Status("MATCHED")

	// StatusInfoSubmitted is returned when the addresses of the traders is
	// submitted.
	StatusInfoSubmitted = Status("INFO_SUBMITTED")

	// StatusInitiateDetailsAcquired is returned when the swapper acquires the
	// initiate details either from generating them or recieving them from the
	// initiating trader.
	StatusInitiateDetailsAcquired = Status("INITIATE_DETAILS_ACQUIRED")

	// StatusInitiated is returned when the atomic swap is initiated.
	StatusInitiated = Status("INITIATED")

	// StatusWaitingForCounterInitiation is returned when the trader is waiting
	// for the counter party to initiate.
	StatusWaitingForCounterInitiation = Status("WAITING_FOR_COUNTER_INITIATION")

	// StatusRedeemDetailsAcquired is returned when the swapper acquires the
	// secret, either from it's local storage or from the blockchain.
	StatusRedeemDetailsAcquired = Status("REDEEM_DETAILS_ACQUIRED")

	// StatusRedeemed is returned when the atomic swap is redeemed.
	StatusRedeemed = Status("REDEEMED")

	// StatusWaitingForCounterRedemption is returned when the swapper is waiting
	// for the counter party initiation.
	StatusWaitingForCounterRedemption = Status("WAITING_FOR_COUNTER_REDEMPTION")

	// StatusRefunded is returned when the given atomic swap is refunded.
	StatusRefunded = Status("REFUNDED")

	// StatusComplained is returned when the swapper complains about a failed
	// atomic swap.
	StatusComplained = Status("COMPLAINED")

	// StatusReceivedSwapDetails is returned when the swapper receives swap
	// details of the counter-party.
	StatusReceivedSwapDetails = Status("RECEIVED_SWAP_DETAILS")

	// StatusSentSwapDetails is returned after the swapper sent it's swap
	// details.
	StatusSentSwapDetails = Status("SENT_SWAP_DETAILS")

	// StatusAudited is returned after the swap is audited.
	StatusAudited = Status("AUDITED")

	// StatusExpired is returned if the given swap is expred.
	StatusExpired = Status("EXPIRED")
)

type Request struct {
	UID                [32]byte    `json:"uid"`
	SendToken          token.Token `json:"sendToken"`
	ReceiveToken       token.Token `json:"receiveToken"`
	SendValue          string      `json:"sendValue"`
	ReceiveValue       string      `json:"receiveAddress"`
	SentToAddress      string      `json:"sentToAddress"`
	ReceiveFromAddress string      `json:"receiveFromAddress"`
	GoesFirst          bool        `json:"goesFirst"`
}
