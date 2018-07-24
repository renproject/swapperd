package store

type Store interface {
	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error
	Delete([]byte) error
}
