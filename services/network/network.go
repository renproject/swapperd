package network

type Network interface {
	Send([32]byte, []byte) error
	Recieve([32]byte) ([]byte, error)
}
