package swap

type SwapAdapter interface {
	SetOwnerAddress([32]byte, []byte) error
	GetOwnerAddress([32]byte) ([]byte, error)
	ReceiveSwapDetails([32]byte) ([]byte, error)
	SendSwapDetails([32]byte, []byte) error
}
