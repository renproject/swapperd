package foundation

type StatusUpdate struct {
	ID     SwapID `json:"id"`
	Status int    `json:"status"`
}

func NewStatusUpdate(id SwapID, status int) StatusUpdate {
	return StatusUpdate{id, status}
}

type StatusQuery struct {
	Responder chan<- map[SwapID]SwapStatus
}
