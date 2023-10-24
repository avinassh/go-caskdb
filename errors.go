package caskdb

import "errors"

var (
	ErrKeyNotExist       = errors.New("invalid key: key doesn't exist")
	ErrKeyNotFound       = errors.New("invalid key: key either deleted or expired")
	ErrKeyValueCorrupted = errors.New("corrupt value: checksum failed")
	ErrEncodingFailed    = errors.New("encoding fail: failed to encode kv record")
	ErrDecodingFailed    = errors.New("decoding fail: failed to decode kv record")
	ErrSeekFailed        = errors.New("see fail: failed to seek to the correct offset")
	ErrReadFailed        = errors.New("read fail: failed to read data from disk")
	ErrInvalidValue      = errors.New("invalid value: trying to store unsupported value")
)
