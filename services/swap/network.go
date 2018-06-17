package swap

type Network interface {
	RecieveSwapDetails([32]byte) ([]byte, error)
	SendSwapDetails([32]byte, []byte) error
}
