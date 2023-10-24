package caskdb

import "time"

type Store interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	SetX(key string, value interface{}, expiry time.Duration) error
	Delete(key string) error
	ListKeys(string) <-chan Record
	Close() bool
}
