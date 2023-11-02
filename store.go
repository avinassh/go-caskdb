package caskdb

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
	ListKeys() []string
	Close() bool
}

var _ Store = (*DiskStore)(nil)
