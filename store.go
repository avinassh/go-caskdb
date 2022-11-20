package caskdb

type Store interface {
	Get(key string) string
	Set(key string, value string)
	Close() bool
}
