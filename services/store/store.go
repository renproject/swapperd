package store

type store struct {
}

type Store interface {
	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error
}
