package swap

const (
	Inactive = iota
	Initiated
	Audited
	AuditPending
	AuditFailed
	Redeemed
	AuditedSecret
	Refunded
	RefundFailed
	Cancelled
	Expired
)

type StatusUpdate struct {
	ID   SwapID `json:"id"`
	Code int    `json:"status"`
}

func NewStatusUpdate(id SwapID, status int) StatusUpdate {
	return StatusUpdate{id, status}
}
