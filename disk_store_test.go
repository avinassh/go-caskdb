package caskdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestDiskStore_Get(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")
	store.Set("name", "jojo")
	val, _ := store.Get("name")
	if val != "jojo" {
		t.Errorf("Get() = %v, want %v", val, "jojo")
	}
}

func TestDiskStore_GetInvalid(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")
	val, _ := store.Get("some key")
	if val != "" {
		t.Errorf("Get() = %v, want %v", val, "")
	}
}

func TestDiskStore_SetWithPersistence(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]string{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"war and peace":        "tolstoy",
		"hamlet":               "shakespeare",
		"othello":              "shakespeare",
		"brave new world":      "huxley",
		"dune":                 "frank herbert",
	}

	for key, val := range tests {
		store.Set(key, val)
		actualVal, _ := store.Get(key)
		if actualVal != val {
			t.Errorf("Get() = %v, want %v", actualVal, val)
		}
	}
	store.Close()
	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	for key, val := range tests {
		actualVal, _ := store.Get(key)
		if actualVal != val {
			t.Errorf("Get() = %v, want %v", actualVal, val)
		}
	}
	store.Close()
}

func TestDiskStore_Delete(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]string{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"war and peace":        "tolstoy",
		"hamlet":               "shakespeare",
		"othello":              "shakespeare",
		"brave new world":      "huxley",
		"dune":                 "frank herbert",
	}
	for key, val := range tests {
		store.Set(key, val)
	}

	// only for tests
	deletedKeys := []string{"hamlet", "dune", "othello"}
	//delete few keys
	for _, k := range deletedKeys {
		store.Delete(k)
	}
	store.Close()

	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}

	//check for deletion
	for _, dkeys := range deletedKeys {
		actualVal, err := store.Get(dkeys)

		if actualVal != "" {
			t.Errorf("Get() = %s, want %s", actualVal, "")
		}

		if errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Get() = %v, want %v", err, ErrKeyNotFound)
		}
	}
	store.Close()
}

func TestDiskStore_ValidCheckSum(t *testing.T) {
	store, _ := NewDiskStore("test.db")
	defer store.Close()
	defer os.Remove("test.db")

	k1, v1 := "👋", "world"
	h1 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k1)), ValueSize: uint32(len(v1)), Meta: 0}
	r1 := Record{Header: h1, Key: k1, Value: v1, RecordSize: headerSize + +h1.KeySize + h1.ValueSize}
	r1.Header.CheckSum = r1.CalculateCheckSum()

	k2, v2 := "", ""
	h2 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k2)), ValueSize: uint32(len(v2)), Meta: 1}
	r2 := Record{Header: h2, Key: k2, Value: v2, RecordSize: headerSize + h2.KeySize + h2.ValueSize}
	r2.Header.CheckSum = r2.CalculateCheckSum()

	k3, v3 := "🔑", ""
	h3 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k3)), ValueSize: uint32(len(v3)), Meta: 0}
	r3 := Record{Header: h3, Key: k3, Value: v3, RecordSize: headerSize + h3.KeySize + h3.ValueSize}
	r3.Header.CheckSum = r3.CalculateCheckSum()

	tests := []Record{r1, r2, r3}

	// valid checksum
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		tt.EncodeKV(buf)

		// store the data
		store.keyDir[tt.Key] = NewKeyEntry(tt.Header.TimeStamp, uint32(store.writePosition), tt.Size())
		store.writePosition += int(tt.Size())
		store.write(buf.Bytes())

		// retrieve the data
		kEntry := store.keyDir[tt.Key]

		//seek to the record
		store.file.Seek(int64(kEntry.position), defaultWhence)

		kvRecord := make([]byte, kEntry.totalSize)
		_, err := io.ReadFull(store.file, kvRecord)
		if err != nil {
			t.Errorf("error in reading the record: %v", err)
		}

		result := &Record{}
		err = result.DecodeKV(kvRecord)
		if err != nil {
			t.Errorf("error in decoding the record: %v", err)
		}

		expectedCheckSum := tt.Header.CheckSum
		actualCheckSum := result.Header.CheckSum

		if expectedCheckSum != actualCheckSum {
			t.Errorf("invalid checksum: Got: %d, Want: %d", actualCheckSum, expectedCheckSum)
		}
	}
}

func TestDiskStore_InValidCheckSum(t *testing.T) {
	store, _ := NewDiskStore("test.db")
	defer store.Close()
	defer os.Remove("test.db")

	k1, v1 := "👋", "world"
	h1 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k1)), ValueSize: uint32(len(v1)), Meta: 0}
	r1 := Record{Header: h1, Key: k1, Value: v1, RecordSize: headerSize + +h1.KeySize + h1.ValueSize}
	r1.Header.CheckSum = r1.CalculateCheckSum()

	k2, v2 := "mykey", ""
	h2 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k2)), ValueSize: uint32(len(v2)), Meta: 1}
	r2 := Record{Header: h2, Key: k2, Value: v2, RecordSize: headerSize + h2.KeySize + h2.ValueSize}
	r2.Header.CheckSum = r2.CalculateCheckSum()

	k3, v3 := "🔑", ""
	h3 := Header{TimeStamp: uint32(time.Now().Unix()), KeySize: uint32(len(k3)), ValueSize: uint32(len(v3)), Meta: 0}
	r3 := Record{Header: h3, Key: k3, Value: v3, RecordSize: headerSize + h3.KeySize + h3.ValueSize}
	r3.Header.CheckSum = r3.CalculateCheckSum()

	tests := []Record{r1, r2, r3}

	// valid checksum
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		tt.EncodeKV(buf)

		// store the data
		store.keyDir[tt.Key] = NewKeyEntry(tt.Header.TimeStamp, uint32(store.writePosition), tt.Size())
		store.writePosition += int(tt.Size())
		store.write(buf.Bytes())

		// retrieve the data
		kEntry := store.keyDir[tt.Key]

		// seek to the record
		store.file.Seek(int64(kEntry.position), defaultWhence)

		kvRecord := make([]byte, kEntry.totalSize)
		_, err := io.ReadFull(store.file, kvRecord)
		if err != nil {
			t.Errorf("error in reading the record: %v", err)
		}

		// corrupt the record by overriding few bytes with corruptedBytes
		corruptedBytes := []byte{12, 90, 87, 101}
		start, end := 17, 21
		for i := start; i < end; i++ {
			kvRecord[i] = corruptedBytes[i-start]
		}

		// write the corrupted bytes and update the hash table
		store.keyDir[tt.Key] = NewKeyEntry(tt.Header.TimeStamp+uint32(time.Now().Unix()), uint32(store.writePosition), tt.Size())
		store.writePosition += int(tt.Size())
		store.write(kvRecord)

		// read the corrupted record
		kEntry, _ = store.keyDir[tt.Key]
		store.file.Seek(int64(kEntry.position), defaultWhence)

		corruptedKV := make([]byte, kEntry.totalSize)
		_, err = io.ReadFull(store.file, corruptedKV)
		if err != nil {
			t.Errorf("error in reading the record: %v", err)
		}

		result := &Record{}
		err = result.DecodeKV(corruptedKV)
		if err != nil {
			t.Errorf("error while decoding the kv record: %v", err)
		}

		expectedCheckSum := tt.Header.CheckSum
		actualCheckSum := result.CalculateCheckSum()

		if expectedCheckSum == actualCheckSum {
			t.Error("checksum matched: data is supposed to be corrupted!")
		}
	}
}

func BenchmarkSetOperation(b *testing.B) {
	store, _ := NewDiskStore("test.db")
	defer store.Close()

	table := map[string]string{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"war and peace":        "tolstoy",
		"hamlet":               "shakespeare",
		"othello":              "shakespeare",
		"brave new world":      "huxley",
		"dune":                 "frank herbert",
	}

	for k, v := range table {
		b.Run(fmt.Sprintf("key:%s", k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				store.Set(k, v)
			}
		})
	}
}

func BenchmarkGetOperation(b *testing.B) {
	store, _ := NewDiskStore("test.db")
	defer store.Close()
	defer os.Remove("test.db")

	table := map[string]string{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"war and peace":        "tolstoy",
		"hamlet":               "shakespeare",
		"othello":              "shakespeare",
		"brave new world":      "huxley",
		"dune":                 "frank herbert",
	}

	for k := range table {
		b.Run(fmt.Sprintf("key:%s", k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				store.Get(k)
			}
		})
	}
}
