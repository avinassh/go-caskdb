package caskdb

import "errors"

var (
	ErrEmptyKey    = errors.New("invalid key: empty key not allowed")
	ErrLargeKey    = errors.New("invalid key: size cant be greater than 4.2GB")
	ErrKeyNotFound = errors.New("invalid key: key either deleted or expired")

	ErrLargeValue = errors.New("invalid value: size cant be greater than 4.2GB")

	ErrSeekFailed = errors.New("see fail: failed to seek to the correct offset")
	ErrReadFailed = errors.New("read fail: failed to read data from disk")

	ErrEncodingFailed = errors.New("encoding fail: failed to encode kv record")
	ErrDecodingFailed = errors.New("decoding fail: failed to decode kv record")
)
