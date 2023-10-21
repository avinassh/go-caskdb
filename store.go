package caskdb

import "time"

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	SetX(key string, value string, expiry time.Duration) error
	Delete(key string) error
	Close() bool
}
