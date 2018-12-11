package swap

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

// StatusUpdate shows the status change of a swap.
type StatusUpdate struct {
	ID     SwapID `json:"id"`
	Status Status `json:"status"`
}

// NewStatusUpdate creates a new `StatusUpdate` with given swap ID and status.
func NewStatusUpdate(id SwapID, status Status) StatusUpdate {
	return StatusUpdate{id, status}
}
