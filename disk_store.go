package caskdb

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

const defaultWhence = 0

type DiskStore struct {
	file          *os.File
	writePosition int
	keyDir        map[string]KeyEntry
}

func NewDiskStore(fileName string) (*DiskStore, error) {
	ds := &DiskStore{keyDir: make(map[string]KeyEntry)}
	if _, err := os.Stat(fileName); err == nil || errors.Is(err, fs.ErrExist) {
		ds.initKeyDir(fileName)
	}
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	ds.file = file
	return ds, nil
}

func (d *DiskStore) Get(key string) string {
	kv, ok := d.keyDir[key]
	if !ok {
		return ""
	}
	// TODO: handle errors
	d.file.Seek(int64(kv.position), defaultWhence)
	data := make([]byte, kv.totalSize)
	// TODO: handle errors
	_, err := io.ReadFull(d.file, data)
	if err != nil {
		panic("read error")
	}
	_, _, value := decodeKV(data)
	return value
}

func (d *DiskStore) Set(key string, value string) {
	timestamp := uint32(time.Now().Unix())
	size, data := encodeKV(timestamp, key, value)
	d.write(data)
	d.keyDir[key] = NewKeyEntry(timestamp, uint32(d.writePosition), uint32(size))
	d.writePosition += size
}

func (d *DiskStore) Close() bool {
	// TODO: handle errors
	d.file.Sync()
	if err := d.file.Close(); err != nil {
		// TODO: log the error
		return false
	}
	return true
}

func (d *DiskStore) write(data []byte) {
	// TODO: handle errors
	_, err := d.file.Write(data)
	if err != nil {
		panic(err)
	}
	// TODO: handle errors
	d.file.Sync()
}

func (d *DiskStore) initKeyDir(existingFile string) {
	file, _ := os.Open(existingFile)
	defer file.Close()
	for {
		header := make([]byte, headerSize)
		_, err := io.ReadFull(file, header)
		if err == io.EOF {
			break
		}
		// TODO: handle errors
		if err != nil {
			break
		}
		timestamp, keySize, valueSize := decodeHeader(header)
		key := make([]byte, keySize)
		value := make([]byte, valueSize)
		_, err = io.ReadFull(file, key)
		// TODO: handle errors
		if err != nil {
			break
		}
		_, err = io.ReadFull(file, value)
		// TODO: handle errors
		if err != nil {
			break
		}
		totalSize := headerSize + keySize + valueSize
		d.keyDir[string(key)] = NewKeyEntry(timestamp, uint32(d.writePosition), totalSize)
		d.writePosition += int(totalSize)
		fmt.Printf("loaded key=%s, value=%s\n", key, value)
	}
}
