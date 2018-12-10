package foundation

// Status of the swap.
type Status int

const (
	Inactive = Status(iota)
	Initiated
	Audited
	AuditFailed
	Redeemed
	Refunded
)

// The SwapStatus contains the swap details and the status.
type SwapStatus struct {
	ID            SwapID `json:"id"`
	SendToken     string `json:"sendToken"`
	ReceiveToken  string `json:"receiveToken"`
	SendAmount    string `json:"sendAmount"`
	ReceiveAmount string `json:"receiveAmount"`
	Timestamp     int64  `json:"timestamp"`
	Status        Status `json:"status"`
}

// StatusUpdate shows the status change of a swap.
type StatusUpdate struct {
	ID     SwapID `json:"id"`
	Status Status `json:"status"`
}

// NewStatusUpdate creates a new `StatusUpdate` with given swap ID and status.
func NewStatusUpdate(id SwapID, status Status) StatusUpdate {
	return StatusUpdate{id, status}
}

// StatusQuery is provided when querying statuses of swaps. It contains a
// Responsder channel which responses can be written into it.
type StatusQuery struct {
	Responder chan<- map[SwapID]SwapStatus
}

