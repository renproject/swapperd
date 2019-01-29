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
