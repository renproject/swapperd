package swap

const (
	StatusUnknown                     = "UNKNOWN"
	StatusMatched                     = "MATCHED"
	StatusInfoSubmitted               = "INFO_SUBMITTED"
	StatusInitateDetailsAcquired      = "INITIATE_DETAILS_ACQUIRED"
	StatusInitiated                   = "INITIATED"
	StatusWaitingForCounterInitiation = "WAITING_FOR_COUNTER_INITIATION"
	StatusRedeemDetailsAcquired       = "REDEEM_DETAILS_ACQUIRED"
	StatusRedeemed                    = "REDEEMED"
	StatusWaitingForCounterRedemption = "WAITING_FOR_COUNTER_REDEMPTION"
	StatusRefunded                    = "REFUNDED"
	StatusComplained                  = "COMPLAINED"

	StatusReceivedSwapDetails = "ReceiveD_SWAP_DETAILS"
	StatusSentSwapDetails     = "SENT_SWAP_DETAILS"
	StatusAudited             = "AUDITED"
)
