package swap

type Network interface {
	ReceiveSwapDetails([32]byte) ([]byte, error)
	SendSwapDetails([32]byte, []byte) error
}
