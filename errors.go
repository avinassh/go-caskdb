package caskdb

import "errors"

var (
	ErrKeyNotFound    = errors.New("invalid key: key either deleted or expired")
	ErrSeekFailed     = errors.New("see fail: failed to seek to the correct offset")
	ErrReadFailed     = errors.New("read fail: failed to read data from disk")
	ErrEncodingFailed = errors.New("encoding fail: failed to encode kv record")
)
