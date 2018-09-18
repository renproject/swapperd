package swap

// ID is the swap ID
type ID [32]byte

type Status string

const (
	StatusUnknown                     = Status("UNKNOWN")
	StatusMatched                     = Status("MATCHED")
	StatusInfoSubmitted               = Status("INFO_SUBMITTED")
	StatusInitiateDetailsAcquired     = Status("INITIATE_DETAILS_ACQUIRED")
	StatusInitiated                   = Status("INITIATED")
	StatusWaitingForCounterInitiation = Status("WAITING_FOR_COUNTER_INITIATION")
	StatusRedeemDetailsAcquired       = Status("REDEEM_DETAILS_ACQUIRED")
	StatusRedeemed                    = Status("REDEEMED")
	StatusWaitingForCounterRedemption = Status("WAITING_FOR_COUNTER_REDEMPTION")
	StatusRefunded                    = Status("REFUNDED")
	StatusComplained                  = Status("COMPLAINED")

	StatusReceivedSwapDetails = Status("ReceiveD_SWAP_DETAILS")
	StatusSentSwapDetails     = Status("SENT_SWAP_DETAILS")
	StatusAudited             = Status("AUDITED")
)
