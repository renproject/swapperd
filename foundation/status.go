package foundation

const (
	Inactive = iota
	Initiated
	Audited
	AuditFailed
	Redeemed
	Refunded
)

// StatusUpdate shows the status change of a swap.
type StatusUpdate struct {
	ID     SwapID `json:"id"`
	Status int    `json:"status"`
}

// NewStatusUpdate creates a new `StatusUpdate` with given ID and status.
func NewStatusUpdate(id SwapID, status int) StatusUpdate {
	return StatusUpdate{id, status}
}

type StatusQuery struct {
	Responder chan<- map[SwapID]SwapStatus
}
