package caskdb

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
	Close() bool
}
