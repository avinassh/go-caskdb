package caskdb

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_encodeHeader(t *testing.T) {
	tests := []*Header{
		{10, 10, 10, 10, 10, IntTypeID},
		{0, 0, 0, 0, 0, UIntTypeID},
		{10000, 10000, 10000, 10000, 10000, StrTypeID},
	}
	for _, tt := range tests {
		newBuf := make([]byte, headerSize)
		//encode the header
		tt.EncodeHeader(newBuf)

		//encoded header should be 24bytes
		if len(newBuf) != headerSize {
			t.Errorf("Invalid encode: expected header size = %v, got = %v", headerSize, len(newBuf))
		}

		//decode the header
		result := &Header{}
		result.DecodeHeader(newBuf)

		if result.CheckSum != tt.CheckSum {
			t.Errorf("EncodeHeader() checksum = %v, want %v", result.CheckSum, tt.CheckSum)
		}
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
	//prepare record
	k1, v1 := "hello", UglyRune('%')
	h1 := Header{TimeStamp: uint32(time.Now().Unix()), ExpiryTime: 0, KeySize: uint32(len(k1))}
	r1 := Record{Header: h1, Key: k1}
	err := r1.EncodeValue(v1)
	if err != nil {
		fmt.Println("Err in encoding the value", err)
		os.Exit(-1)
	}
	r1.Header.ValueSize = uint32(len(r1.Value))
	r1.Header.CheckSum = r1.Header.CalculateCheckSum(r1.Value)

	k2, v2 := "", 0.000121289323
	h2 := Header{TimeStamp: uint32(time.Now().Unix()), ExpiryTime: uint32(2 * time.Second), KeySize: uint32(len(k2))}
	r2 := Record{Header: h2, Key: k2}
	err = r2.EncodeValue(v2)
	if err != nil {
		fmt.Println("Err in encoding the value", err)
		os.Exit(-1)
	}
	r2.Header.ValueSize = uint32(len(r2.Value))
	r2.Header.CheckSum = r2.Header.CalculateCheckSum(r2.Value)

	k3, v3 := "🔑", -12.2901
	h3 := Header{TimeStamp: uint32(time.Now().Unix()), ExpiryTime: 0, KeySize: uint32(len(k3))}
	r3 := Record{Header: h3, Key: k3}
	err = r3.EncodeValue(v3)
	if err != nil {
		fmt.Println("Err in encoding the value", err)
		os.Exit(-1)
	}
	r3.Header.ValueSize = uint32(len(r3.Value))
	r3.Header.CheckSum = r3.Header.CalculateCheckSum(r3.Value)

	tests := []Record{r1, r2, r3}
	for _, tt := range tests {
		//encode the record
		buf := bytes.NewBuffer(make([]byte, headerSize))
		tt.EncodeKV(buf)

		//encoded buffer size should be equal to headersize + keysize + valuesize
		expectedSize := (headerSize + tt.Header.KeySize + tt.Header.ValueSize)
		if uint32(len(buf.Bytes())) != expectedSize {
			t.Errorf("EncodeKV() invalid encoding, expected size=%v, got=%v", expectedSize, uint32(len(buf.Bytes())))
		}

		//decode the record
		result := &Record{}
		result.DecodeKV(buf.Bytes())

		if result.Header.CalculateCheckSum(result.Value) != tt.Header.CheckSum {
			t.Errorf("EncodeKV() checksum = %v, want %v", result.Header.CalculateCheckSum(result.Value), tt.Header.CheckSum)
		}
		if result.Header.TimeStamp != tt.Header.TimeStamp {
			t.Errorf("EncodeKV() timestamp = %v, want %v", result.Header.TimeStamp, tt.Header.TimeStamp)
		}
		if result.Key != tt.Key {
			t.Errorf("EncodeKV() key = %v, want %v", result.Key, tt.Key)
		}

		actualVal, err := result.DecodeValue()
		if err != nil {
			fmt.Println("error while decoding actual value", err)
			os.Exit(-1)
		}
		expVal, err := tt.DecodeValue()
		if err != nil {
			fmt.Println("error while decoding expected value", err)
			os.Exit(-1)
		}
		if actualVal != expVal {
			t.Errorf("encodeKV() value = %v, want %v", result.Value, tt.Value)
		}

	}
}
