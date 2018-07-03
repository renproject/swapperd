package swap

type Info interface {
	SetOwnerAddress([32]byte, []byte) error
	GetOwnerAddress([32]byte) ([]byte, error)
}
