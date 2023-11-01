package caskdb

import (
	"bytes"
	"testing"
	"time"
)

func Test_encodeHeader(t *testing.T) {
	tests := []*Header{
		{10, 10, 10, 1},
		{0, 0, 0, 0},
		{10000, 10000, 10000, 1},
	}
	for _, tt := range tests {
		newBuf := new(bytes.Buffer)
		//encode the header
		tt.EncodeHeader(newBuf)

		//encoded header should be 14bytes
		if len(newBuf.Bytes()) != headerSize {
			t.Errorf("Invalid encode: expected header size = %v, got = %v", headerSize, len(newBuf.Bytes()))
		}

		//decode the header
		result := &Header{}
		result.DecodeHeader(newBuf.Bytes())

		if result.TimeStamp != tt.TimeStamp {
			t.Errorf("EncodeHeader() timestamp = %v, want %v", result.TimeStamp, tt.TimeStamp)
		}
		if result.KeySize != tt.KeySize {
			t.Errorf("EncodeHeader() keySize = %v, want %v", result.KeySize, tt.KeySize)
		}
		if result.ValueSize != tt.ValueSize {
			t.Errorf("EncodeHeader() valueSize = %v, want %v", result.ValueSize, tt.ValueSize)
		}
	}
}

func Test_encodeKV(t *testing.T) {
	k1, v1 := "hello", "world"
	h1 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k1)), ValueSize: uint32(len(v1)), Meta: 0}
	r1 := Record{Header: h1, Key: k1, Value: v1, RecordSize: headerSize + +h1.KeySize + h1.ValueSize}

	k2, v2 := "", ""
	h2 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k2)), ValueSize: uint32(len(v2)), Meta: 1}
	r2 := Record{Header: h2, Key: k2, Value: v2, RecordSize: headerSize + h2.KeySize + h2.ValueSize}

	k3, v3 := "🔑", ""
	h3 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k3)), ValueSize: uint32(len(v3)), Meta: 0}
	r3 := Record{Header: h3, Key: k3, Value: v3, RecordSize: headerSize + h3.KeySize + h3.ValueSize}

	tests := []Record{r1, r2, r3}
	for _, tt := range tests {
		//encode the record
		buf := new(bytes.Buffer)
		tt.EncodeKV(buf)

		//encoded buffer size should be equal to headersize + keysize + valuesize
		expectedSize := tt.RecordSize
		if uint32(len(buf.Bytes())) != expectedSize {
			t.Errorf("EncodeKV() invalid encoding, expected size=%v, got=%v", expectedSize, uint32(len(buf.Bytes())))
		}

		//decode the record
		result := &Record{}
		result.DecodeKV(buf.Bytes())

		if result.Header.TimeStamp != tt.Header.TimeStamp {
			t.Errorf("EncodeKV() timestamp = %v, want %v", result.Header.TimeStamp, tt.Header.TimeStamp)
		}
		if result.Key != tt.Key {
			t.Errorf("EncodeKV() key = %v, want %v", result.Key, tt.Key)
		}

		if result.Value != tt.Value {
			t.Errorf("encodeKV() value = %v, want %v", result.Value, tt.Value)
		}

	}
}
