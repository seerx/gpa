package dbutil

type BlobReadWriter interface {
	Read([]byte) error
	Write() ([]byte, error)
}
